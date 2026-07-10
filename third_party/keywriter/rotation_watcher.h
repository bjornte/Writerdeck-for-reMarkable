#ifndef ROTATION_WATCHER_H
#define ROTATION_WATCHER_H

#include <QObject>

// Relays QML rotationChanged to Writerdeck-server via a C callback.
class RotationWatcher : public QObject
{
    Q_OBJECT
public:
    void setRoot(QObject *root) { m_root = root; }
    void setApplying(bool *flag) { m_applyingFlag = flag; }
    void setNotify(void (*fn)(int)) { m_notify = fn; }

public slots:
    void onRotationChanged();

private:
    QObject *m_root = nullptr;
    bool *m_applyingFlag = nullptr;
    void (*m_notify)(int) = nullptr;
};

#endif
