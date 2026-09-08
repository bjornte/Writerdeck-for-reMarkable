#include <QDir>
#include <QGuiApplication>
#include <QQmlApplicationEngine>
#include <QQuickWindow>
#include <QTimer>

int main(int argc, char *argv[]) {
    QGuiApplication app(argc, argv);

    QQmlApplicationEngine engine;
    QObject::connect(
        &engine, &QQmlApplicationEngine::objectCreationFailed, &app,
        []() { QCoreApplication::exit(-1); }, Qt::QueuedConnection);
    engine.loadFromModule("remarkable_example", "Main");

    // Spike: PNG for scripts/capture-screenshot.sh (grabToImage fails on epaper).
    const QString shotPath = QDir::currentPath() + QStringLiteral("/screen.png");
    QTimer::singleShot(3000, &app, [&engine, shotPath]() {
        const QObjectList roots = engine.rootObjects();
        if (roots.isEmpty())
            return;
        auto *win = qobject_cast<QQuickWindow *>(roots.first());
        if (!win)
            return;
        win->grabWindow().save(shotPath, "PNG");
    });

    return app.exec();
}
