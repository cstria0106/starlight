import os
import subprocess
import argparse

def run_and_wait(image: str, cmd: str, wait_for: str, args: str):
    if args is None:
        args = ' '

    os.system('sudo ctr-starlight pull --profile myproxy cloud.cluster.local/%s' % image)
    os.system('sudo ctr c create \
        --snapshotter=starlight \
        %s \
        cloud.cluster.local/%s \
        instance1 %s' % (args, image, cmd))
        
    p = subprocess.Popen('sudo ctr task start instance1', shell=True, stderr=subprocess.STDOUT, stdout=subprocess.PIPE)
    while True:
        l = p.stdout.readline()
        if l == '':
            continue
        
        print('[stdout] %s' % l.strip())
        if wait_for is not None:
            if str(l).find(wait_for) >= 0:
                print('[done]')
                rc = os.system('sudo ctr task kill instance1')
                assert(rc == 0)
                break

    p.wait()

def main():
    parser = argparse.ArgumentParser(description='Start up time benchmark tool for Starlight')
    parser.add_argument('image', type=str)
    parser.add_argument('cmd', type=str)
    parser.add_argument('--wait-for', dest='wait_for', type=str)
    parser.add_argument('--additional-args', dest='args', type=str)
    args = parser.parse_args()
    run_and_wait(args.image, args.cmd, args.wait_for, args.args)
    pass

if __name__ == '__main__':
    main()
