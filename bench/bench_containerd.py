
import argparse
import os
from typing import Iterable, Tuple

from .bench import Service, ShellCommand, StartTimerCommand, StopTimerCommand, TimerContext, MarkTimerCommand


class _ContainerdService(Service):
    __mounts: Iterable[Tuple[str, str]]

    def __init__(self, image: str, cmd: str, wait_for: str, env: dict[str, str], mounts: Iterable[Tuple[str, str]], timer_index: int, output: str | None) -> None:
        start_args = '--insecure-registry --name instance '

        for key, value in env.items():
            start_args += '--env %s=%s ' % (key, value)

        self.__mounts = mounts
        for src, dst in mounts:
            start_args += '--mount type=bind,source=%s,target=%s' % (src, dst)

        start_command = 'sudo nerdctl run %s %s %s' % (start_args, image, cmd)

        timer_context = TimerContext(timer_index, output=output)

        super().__init__((
            StartTimerCommand(timer_context),
            ShellCommand(start_command, wait_for,
                         (MarkTimerCommand(timer_context), ShellCommand('sudo nerdctl container kill -s INT instance'))),
            ShellCommand('sudo nerdctl container rm instance'),
            StopTimerCommand(timer_context)
        ))

    def run(self):
        for src, _ in self.__mounts:
            os.makedirs(src, exist_ok=True)
        return super().run()


class __ContainerdServiceBuilder:
    image: str
    cmd: str
    wait_for: str
    env: dict[str, str]
    mounts: Iterable[Tuple[str, str]]

    def __init__(self, image: str, cmd: str, wait_for: str, env: dict[str, str], mounts: Iterable[Tuple[str, str]]) -> None:
        self.image = image
        self.cmd = cmd
        self.wait_for = wait_for
        self.env = env
        self.mounts = mounts

    def build(self, timer_index: int, output_name: str | None) -> _ContainerdService:
        return _ContainerdService(self.image, self.cmd, self.wait_for, self.env, self.mounts, timer_index, output=output_name)


__SERVICE_BUILDERS: dict[str, __ContainerdServiceBuilder] = {
    'redis': __ContainerdServiceBuilder(
        'cloud.cluster.local/redis:6.2.1',
        '/usr/local/bin/redis-server',
        'Ready to accept connections',
        dict(),
        [('/tmp/test-redis-data', '/data')]
    )
}


class Arguments:
    service: str
    timer_index: int
    output: str


def run(args: Arguments):
    if args.service not in __SERVICE_BUILDERS:
        print('No service named \'%s\'' % args.service)
        exit(1)

    exit(__SERVICE_BUILDERS[args.service].build(
        args.timer_index, args.output).run())


def init_argument_parser(parser: argparse.ArgumentParser) -> argparse.ArgumentParser:
    parser.description = 'Check start up time of Containerd container using Nerdctl'
    parser.add_argument('service', type=str, help='service name')
    parser.add_argument('-o', type=str,
                        dest='output', help='path of timer output directory', default=None)
    parser.add_argument('-i', type=str, dest='timer_index',
                        help='index of timer context')
    return parser


def __main():
    args = init_argument_parser(argparse.ArgumentParser()).parse_args()
    run(args)


if __name__ == '__main__':
    __main()
