#include "socket_feed.h"

#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQmlContext>
#include <QQuickWindow>
#include <QTimer>

int main(int argc, char *argv[]) {
    QGuiApplication app(argc, argv);

    const char *sockPath = "/home/root/spike-rm2-socket/Writerdeck.sock";
    const char *shotPath = "/home/root/spike-rm2-socket/screen.png";

    SocketFeed feed;
    if (!feed.listen(QString::fromUtf8(sockPath)))
        return 2;

    QQmlApplicationEngine engine;
    engine.rootContext()->setContextProperty(QStringLiteral("socketFeed"), &feed);
    QObject::connect(
        &engine, &QQmlApplicationEngine::objectCreationFailed, &app,
        []() { QCoreApplication::exit(-1); }, Qt::QueuedConnection);
    engine.loadFromModule("remarkable_socket_spike", "Main");

    QTimer::singleShot(5000, &app, [&engine, shotPath]() {
        const QObjectList roots = engine.rootObjects();
        if (roots.isEmpty())
            return;
        auto *win = qobject_cast<QQuickWindow *>(roots.first());
        if (!win)
            return;
        win->grabWindow().save(QString::fromUtf8(shotPath), "PNG");
    });

    return app.exec();
}
