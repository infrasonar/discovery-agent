#!/usr/bin/env python
import os
import re
import subprocess

platforms = [
    ["solaris", "amd64"],
    ["darwin", "amd64"],
    ["linux", "amd64"],
    ["windows", "amd64"],
    # ["linux", "arm64"],
]

targetpath = 'target'
name = 'discovery-agent'


def read_version() -> str:
    with open('version.go', 'r') as fp:
        content = fp.read()
    re_version = re.compile(r'Version\s=\s"([0-9a-z\.\-]+)"')
    m = re_version.findall(content)
    return m[0]


def build(goos: str, goarch: str, name: str, version: str):
    cwd = os.path.dirname(os.path.abspath(__file__))
    print(cwd)

    outfile = f'{name}.exe' \
        if goos == 'windows' else name

    tmp_env = os.environ.copy()
    tmp_env["GOOS"] = goos
    tmp_env["GOARCH"] = goarch
    tmp_env["CGO_ENABLED"] = "0"

    with subprocess.Popen(
            ['go', 'build', '-trimpath', '-o', outfile],
            env=tmp_env,
            cwd=cwd,
            stdout=subprocess.PIPE) as _proc:
        print(f'Building {goos}/{goarch}...')

    cmd = ['zip', '-r'] if goos == 'windows' else ['tar', '-zcf']
    ext = 'zip' if goos == 'windows' else 'tar.gz'
    target = f'{name}-{goos}-{goarch}-{version}.{ext}'
    target = os.path.join(targetpath, target)
    cmd.extend([target, outfile])
    subprocess.call(cmd, cwd=cwd)
    os.unlink(outfile)


if __name__ == '__main__':
    version = read_version()

    for goos, goarch in platforms:
        build(goos, goarch, name, version)
