#!/usr/bin/python3

import sys

class RevisionLog(object):
    revision = None
    author = None
    timestamp = None
    info = None
    changes = []
    log = None

    def is_branching(self):
        changes = self.changes
        if len(changes) != 1:
            return False

        c = changes[0]
        if c.type != 'A':
            return False
        if c.from_revision is None:
            return False

        return True

    def get_branching_info(self):
        if not self.is_branching():
            return None
        return self.changes[0]

class Change(object):
    type = None
    path = None
    from_branch = None
    from_revision = None

    def __str__(self):
        if self.from_branch is None:
            return '{} {}'.format(self.type, self.path)
        return '{} {} (from {}:{})'.format(self.type, self.path, self.from_branch, self.from_revision)

    def __repr__(self):
        return self.__str__()

class RevisionLogParser(object):
    def __init__(self, inp):
        self.inp = inp

    def get_revision_logs(self):
        for lines in self._split_logs():
            revlog = self._parse_revision_log(lines)
            yield revlog

    def _split_logs(self):
        def is_empty(lines):
            return len(''.join(lines).strip()) == 0

        lines = []
        index = 0
        for line in self.inp:
            line = line.rstrip()

            if self._is_separator(line):
                if not is_empty(lines):
                    yield lines

                lines = []
                index = 0
                continue

            lines.append(line)

        if not is_empty(lines):
            yield lines

    def _is_separator(self, line):
        return len(line) > 0 and len(line.strip('-')) == 0

    def _parse_changes(self, line):
        c = Change()
        line = line.strip()
        c.type = line[0]
        c.path = line[2:]

        from_branch, from_revision = None, None
        t = ' (from'
        path = c.path
        if t in path:
            pos = path.rfind(t)
            pos2 = path.rfind(':')
            c.from_branch = path[pos+len(t)+1:pos2]
            c.from_revision = int(path[pos2+1:-1])
            c.path = path[:pos]

        return c

    def _parse_revision_log(self, lines):
        revlog = RevisionLog()

        info = lines[0]
        p = [t.strip() for t in info.split('|')]

        revlog.revision = int(p[0][1:])
        revlog.author = p[1]
        revlog.timestamp = p[2]
        revlog.info = p[3]

        assert lines[1].startswith('Changed paths:')

        iterator = lines[2:].__iter__()

        changes = []
        for line in iterator:
            if len(line.strip()) == 0:
                break
            changes.append(self._parse_changes(line))
        revlog.changes = changes

        log = []
        for line in iterator:
            log.append(line)
        revlog.log = '\n'.join(log).rstrip()

        return revlog

def analyze(inp):
    p = RevisionLogParser(inp)
    for revlog in p.get_revision_logs():
        branching = revlog.get_branching_info()
        if branching is not None:
            print(branching)


def main():
    if len(sys.argv) > 1:
        f = open(sys.argv[1])
    else:
        f = sys.stdin

    analyze(f)


if __name__ == '__main__':
    main()

