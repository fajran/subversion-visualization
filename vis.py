#!/usr/bin/python3

import sys

def is_separator(line):
    return len(line) > 0 and len(line.strip('-')) == 0

def split_revision_log(inp):
    def is_empty(lines):
        return len(''.join(lines).strip()) == 0

    lines = []
    index = 0
    for line in inp:
        line = line.rstrip()

        if is_separator(line):
            if not is_empty(lines):
                yield lines

            lines = []
            index = 0
            continue

        lines.append(line)

    if not is_empty(lines):
        yield lines

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

def parse_revision_log(lines):
    revlog = RevisionLog()

    info = lines[0]
    p = [t.strip() for t in info.split('|')]

    revlog.revision = int(p[0][1:])
    revlog.author = p[1]
    revlog.timestamp = p[2]
    revlog.info = p[3]

    def parse_changes(line):
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

    assert lines[1].startswith('Changed paths:')

    iterator = lines[2:].__iter__()

    changes = []
    for line in iterator:
        if len(line.strip()) == 0:
            break
        changes.append(parse_changes(line))
    revlog.changes = changes

    log = []
    for line in iterator:
        log.append(line)
    revlog.log = '\n'.join(log).rstrip()

    return revlog

def analyze(inp):
    for log in split_revision_log(inp):
        revlog = parse_revision_log(log)

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

