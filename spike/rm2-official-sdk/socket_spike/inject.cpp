#include <cstdio>
#include <cstring>
#include <unistd.h>
#include <sys/socket.h>
#include <sys/un.h>

static int send_line(int fd, const char *line) {
    const size_t n = strlen(line);
    ssize_t w = write(fd, line, n);
    return (w == static_cast<ssize_t>(n)) ? 0 : 1;
}

int main(int argc, char **argv) {
    const char *path = (argc > 1) ? argv[1] : "/home/root/spike-rm2-socket/Writerdeck.sock";
    int fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (fd < 0)
        return 1;

    sockaddr_un addr {};
    addr.sun_family = AF_UNIX;
    if (strlen(path) >= sizeof(addr.sun_path))
        return 1;
    strncpy(addr.sun_path, path, sizeof(addr.sun_path) - 1);

    if (connect(fd, reinterpret_cast<sockaddr *>(&addr), sizeof(addr)) != 0) {
        close(fd);
        return 1;
    }

    for (const char *p = "SPIKE OK"; *p; ++p) {
        char buf[64];
        snprintf(buf, sizeof(buf), "{\"t\":\"text\",\"cp\":%u}\n", static_cast<unsigned char>(*p));
        if (send_line(fd, buf) != 0) {
            close(fd);
            return 1;
        }
    }
    if (send_line(fd, "{\"t\":\"key\",\"k\":\"Return\"}\n") != 0) {
        close(fd);
        return 1;
    }
    close(fd);
    return 0;
}
