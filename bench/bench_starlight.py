import argparse
import os
from typing import Iterable, Tuple

from .bench import Service, SleepCommand, StartTimerCommand, MarkTimerCommand, StopTimerCommand, TimerContext, ShellCommand


class _StarlightService(Service):
    __mounts: Iterable[Tuple[str, str]]

    def __init__(self, proxy: str, image: str, cmd: str, wait_for: str, env: dict[str, str], mounts: Iterable[Tuple[str, str]], timer_index: int, output: str | None) -> None:
        container_creation_args = ''

        for key, value in env.items():
            container_creation_args += '--env %s=%s ' % (key, value)

        self.__mounts = mounts
        for src, dst in mounts:
            container_creation_args += '--mount type=bind,src=%s,dst=%s,options=rbind:rw ' % (
                src, dst)

        container_creation_args += '--net-host'

        container_creation_cmd = 'sudo ctr containers create --snapshotter=starlight %s %s instance %s' % (
            container_creation_args, image, cmd)

        timer_context = TimerContext(timer_index, output=output)

        super().__init__(
            [
                StartTimerCommand(timer_context),
                ShellCommand(
                    'sudo ctr-starlight pull --profile %s %s' % (proxy, image)),
                ShellCommand(container_creation_cmd),
                ShellCommand('sudo ctr task start instance',
                             wait_for, [MarkTimerCommand(timer_context), ShellCommand('sudo ctr task kill instance')]),
                SleepCommand(5),
                ShellCommand('sudo ctr container rm instance'),
                StopTimerCommand(timer_context)
            ]
        )

    def run(self):
        for src, _ in self.__mounts:
            os.makedirs(src, exist_ok=True)
        return super().run()


class __StarlightServiceBuilder:
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

    def build(self, proxy: str, timer_index: int, output_name: str | None) -> _StarlightService:
        return _StarlightService(proxy, self.image, self.cmd, self.wait_for, self.env, self.mounts, timer_index, output=output_name)


__SERVICE_BUILDERS: dict[str, __StarlightServiceBuilder] = {
    'redis': __StarlightServiceBuilder(
        'cloud.cluster.local/redis:6.2.1-starlight',
        '/usr/local/bin/redis-server',
        'Ready to accept connections',
        dict(),
        [('/tmp/test-redis-data', '/data')]
    )
}


class Arguments:
    service: str
    proxy: str
    timer_index: int
    output: str | None


def run(args: Arguments):
    if args.service not in __SERVICE_BUILDERS:
        print('No service named \'%s\'' % args.service)
        exit(1)

    exit(__SERVICE_BUILDERS[args.service].build(
        args.proxy, args.timer_index, args.output).run())


def init_argument_parser(parser: argparse.ArgumentParser) -> argparse.ArgumentParser:
    parser.description = 'Check start up time of Starlight container'
    parser.add_argument('service', type=str, help='service name')
    parser.add_argument('--proxy', type=str, default='myproxy',
                        help='starlight proxy profile name')
    parser.add_argument('-o', type=str,
                        dest='output', help='path of timer output directory', default=None)
    parser.add_argument('-i', type=int, dest='timer_index',
                        help='index of timer context')
    return parser


def __main():
    args = init_argument_parser(argparse.ArgumentParser()).parse_args()
    run(args)


if __name__ == '__main__':
    __main()
