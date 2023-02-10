
import argparse
import os
from typing import Iterable, Tuple

from .bench import PrintTimerCommand, Service, ShellCommand, StartTimerCommand, TimerContext


class __ContainerdService(Service):
    __mounts: Iterable[Tuple[str, str]]

    def __init__(self, image: str, cmd: str, wait_for: str, env: dict[str, str], mounts: Iterable[Tuple[str, str]]) -> None:
        start_args = '--insecure-registry --name instance '

        for key, value in env.items():
            start_args += '--env %s=%s ' % (key, value)

        self.__mounts = mounts
        for src, dst in mounts:
            start_args += '--mount type=bind,source=%s,target=%s' % (src, dst)

        start_command = 'sudo nerdctl run %s %s %s' % (start_args, image, cmd)

        timer_context = TimerContext('containerd')
        super().__init__((
            StartTimerCommand(timer_context),
            ShellCommand(start_command, wait_for,
                         (PrintTimerCommand(timer_context), ShellCommand('sudo nerdctl container kill -s INT instance'))),
            ShellCommand('sudo nerdctl container rm instance'),
        ))

    def run(self):
        for src, _ in self.__mounts:
            os.makedirs(src, exist_ok=True)
        return super().run()


__SERVICES: dict[str, __ContainerdService] = {
    'redis': __ContainerdService(
        'cloud.cluster.local/redis:6.2.1',
        '/usr/local/bin/redis-server',
        'Ready to accept connections',
        dict(),
        [('/tmp/test-redis-data', '/data')]
    )
}


class Arguments:
    service: str


def run(args: Arguments):
    service_name = args.service

    if service_name not in __SERVICES:
        print('No service named \'%s\'' % service_name)
        exit(1)

    exit(__SERVICES[service_name].run())


def init_argument_parser(parser: argparse.ArgumentParser) -> argparse.ArgumentParser:
    parser.description = 'Check start up time of Containerd container using Nerdctl'
    parser.add_argument('service', type=str)
    return parser


def __main():
    args = init_argument_parser(argparse.ArgumentParser()).parse_args()
    run(args)


if __name__ == '__main__':
    __main()
