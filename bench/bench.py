import argparse
from io import TextIOWrapper
import subprocess
import time
from collections.abc import Iterable


class Command:
    def execute(self):
        raise NotImplementedError()


class TimerContext:
    name: str
    start_time: float | None
    save_as_file: bool
    file: TextIOWrapper | None

    def __init__(self, name: str, save_as_file: bool = False) -> None:
        self.name = name
        self.start_time = None
        self.save_as_file = save_as_file

    def start(self):
        self.start_time = time.time()
        if self.save_as_file:
            self.file = open('timer-%s-%d' %
                             (self.name, int(self.start_time)), 'w')

    def elapsed(self):
        return time.time() - self.start_time

    def mark(self):
        elapsed = self.elapsed()
        print('[timer - %s] %.4fs' %
              (self.name, elapsed))

        if self.file is not None:
            self.file.write('%f,%f' % (self.start_time, elapsed))

    def stop(self):
        if self.file is not None:
            self.file.close()


class StartTimerCommand(Command):
    context: TimerContext

    def __init__(self, context: TimerContext) -> None:
        super().__init__()
        self.context = context

    def execute(self):
        self.context.start()


class MarkTimerCommand(Command):
    context: TimerContext

    def __init__(self, context: TimerContext) -> None:
        super().__init__()
        self.context = context

    def execute(self):
        self.context.mark()


class StopTimerCommand(Command):
    context: TimerContext

    def __init__(self, context: TimerContext) -> None:
        super().__init__()
        self.context = context

    def execute(self):
        self.context.stop()


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
                'shell command \'%s\' returned code %d' % (self.cmd, returncode))


class Service:
    __commands: Iterable[Command]

    def __init__(self, commands: Iterable[Command]) -> None:
        self.__commands = commands

    def run(self):
        for command in self.__commands:
            command.execute()
