import os
import subprocess
import argparse
import time


class Command:
    cmd: str
    wait_for: str | None

    def __init__(self, cmd: str, wait_for: str | None = None) -> None:
        self.cmd = cmd
        self.wait_for = wait_for


class Service:
    __commands: list[Command]

    def __init__(self, commands: list[Command]) -> None:
        self.__commands = commands

    def __execute_command(self, command: Command) -> int:
        p = subprocess.Popen(
            command.cmd, shell=True, stderr=subprocess.STDOUT, stdout=subprocess.PIPE)

        while True:
            returncode = p.poll()
            if returncode is not None:
                return returncode

            l = p.stdout.readline().decode('utf-8')
            if l == '':
                continue

            print('[stdout] %s' % l.strip())
            if command.wait_for is not None:
                if l.find(command.wait_for) >= 0:
                    break

        return p.wait()

    def run(self):
        for command in self.__commands:
            returncode = self.__execute_command(command)
            if returncode != 0:
                print('command \'%s\' has returned %d', command.cmd, returncode)


class StarlightService(Service):
    __mounts: list[(str, str)]

    def __init__(self, image: str, cmd: str, wait_for: str, env: dict[str, str], mounts: list[(str, str)]) -> None:
        container_creation_args = ''

        for key, value in env.items():
            container_creation_args += '--env %s=%s ' % (key, value)

        self.__mounts = mounts
        for src, dst in mounts:
            container_creation_args += '--mount type=bind,src=%s,dst=%s=rbind:rw ' % (
                src, dst)

        container_creation_args += '--net-host'

        container_creation_cmd = 'sudo ctr containers create \
                                    --snapshotter=starlight \
                                    %s \
                                    cloud.cluster.local/%s \
                                    instance %s' % (container_creation_args, image, cmd)

        super().__init__(
            [
                Command(
                    'sudo ctr-starlight pull --profile myproxy cloud.cluster.local/%s' % image),
                Command(container_creation_cmd),
                Command('sudo ctr task start instance', wait_for),
                Command('sudo ctr task kill instance')
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

    SERVICES[service_name].run()


if __name__ == '__main__':
    main()
