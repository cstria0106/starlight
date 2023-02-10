
import argparse
from .bench_starlight import \
    init_argument_parser as init_starlight_argument_parser, \
    run as run_starlight
from .bench_containerd import \
    init_argument_parser as init_containerd_argument_parser, \
    run as run_containerd


def main() -> None:
    parser = argparse.ArgumentParser(
        prog='bench',
        description='Check start up time of Starlight, Containerd container')
    subparsers = parser.add_subparsers()

    starlight_parser = subparsers.add_parser('starlight')
    starlight_parser.set_defaults(func=run_starlight)
    init_starlight_argument_parser(starlight_parser)

    containerd_parser = subparsers.add_parser('containerd')
    containerd_parser.set_defaults(func=run_containerd)
    init_containerd_argument_parser(containerd_parser)

    args = parser.parse_args()
    if 'func' in args:
        args.func(args)
    else:
        parser.print_help()


if __name__ == '__main__':
    main()
