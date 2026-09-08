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
    engine.loadFromModule("remarkable_textedit_spike", "Main");

    const char *shotPath = "/home/root/spike-rm2-textedit/screen.png";
    QTimer::singleShot(3500, &app, [&engine, shotPath]() {
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
