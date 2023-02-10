import os
import subprocess
import argparse
import time
from collections.abc import Iterable


class Command:
    def execute(self):
        raise NotImplementedError()


class TimerContext:
    name: str
    start_time: float | None

    def __init__(self, name: str) -> None:
        self.name = name
        self.start_time = None


class StartTimerCommand(Command):
    context: TimerContext

    def __init__(self, context: TimerContext) -> None:
        super().__init__()
        self.context = context

    def execute(self):
        self.context.start_time = time.time()


class PrintTimerCommand(Command):
    context: TimerContext

    def __init__(self, context: TimerContext) -> None:
        super().__init__()
        self.context = context

    def execute(self):
        now = time.time()
        print('[timer - %s] %.4fs' %
              (self.context.name, now - self.context.start_time))


class SleepCommand(Command):
    __amount: float

    def __init__(self, amount: float) -> None:
        self.__amount = amount

    def execute(self):
        time.sleep(self.__amount)


class ShellCommand(Command):
    cmd: str
    wait_for: str | None
    cleanup_commands: Iterable[Command] | None

    def __init__(self, cmd: str, wait_for: str | None = None, cleanup_commands: Iterable[Command] | None = None) -> None:
        self.cmd = cmd
        self.wait_for = wait_for
        self.cleanup_commands = cleanup_commands
        assert ((self.wait_for is None) == (self.cleanup_commands is None))

    def execute(self):
        print('[run] %s' % self.cmd)
        if self.wait_for is not None:
            print('[wait for] %s' % self.wait_for)
        p = subprocess.Popen(
            self.cmd, shell=True, stderr=subprocess.STDOUT, stdout=subprocess.PIPE)

        returncode = None
        while True:
            returncode = p.poll()
            if returncode is not None:
                return returncode

            l = p.stdout.readline().decode('utf-8')
            if l == '':
                continue

            print('[stdout] %s' % l.strip())
            if self.wait_for is not None:
                if l.find(self.wait_for) >= 0:
                    for cleanup_command in self.cleanup_commands:
                        cleanup_command.execute()
                    break

        if returncode is None:
            returncode = p.wait()

        if returncode != 0:
            raise Exception(
                'shell command \'%s\' returned code %d', self.cmd, returncode)


class Service:
    __commands: Iterable[Command]

    def __init__(self, commands: Iterable[Command]) -> None:
        self.__commands = commands

    def run(self) -> int:
        for command in self.__commands:
            command.execute()

        return 0


class StarlightService(Service):
    __mounts: Iterable[(str, str)]
    __timer_context: TimerContext

    def __init__(self, image: str, cmd: str, wait_for: str, env: dict[str, str], mounts: Iterable[(str, str)]) -> None:
        container_creation_args = ''

        for key, value in env.items():
            container_creation_args += '--env %s=%s ' % (key, value)

        self.__mounts = mounts
        for src, dst in mounts:
            container_creation_args += '--mount type=bind,src=%s,dst=%s,options=rbind:rw ' % (
                src, dst)

        container_creation_args += '--net-host'

        container_creation_cmd = 'sudo ctr containers create \
                                    --snapshotter=starlight \
                                    %s \
                                    cloud.cluster.local/%s \
                                    instance %s' % (container_creation_args, image, cmd)

        self.__timer_context = TimerContext('starlight')

        super().__init__(
            [
                StartTimerCommand(self.__timer_context),
                ShellCommand(
                    'sudo ctr-starlight pull --profile myproxy cloud.cluster.local/%s' % image),
                ShellCommand(container_creation_cmd),
                ShellCommand('sudo ctr task start instance',
                             wait_for, [PrintTimerCommand(self.__timer_context), ShellCommand('sudo ctr task kill instance')]),
                SleepCommand(5),
                ShellCommand('sudo ctr container rm instance')
            ]
        )

    def run(self):
        for src, _ in self.__mounts:
            os.makedirs(src, exist_ok=True)
        return super().run()


SERVICES: dict[str, Service] = {
    'redis': StarlightService(
        'redis:6.2.1-starlight',
        '/usr/local/bin/redis-server',
        'Ready to accept connections',
        dict(),
        [('/tmp/test-redis-data', '/data')]
    )
}


def main():
    parser = argparse.ArgumentParser(
        description='Start up time benchmark tool for Starlight')
    parser.add_argument('service', type=str)
    args = parser.parse_args()

    service_name = args.service
    if service_name not in SERVICES:
        print('No service named \'%s\'' % service_name)
        exit(1)

    exit(SERVICES[service_name].run())


if __name__ == '__main__':
    main()
