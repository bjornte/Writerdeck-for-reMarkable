#ifndef LOBBY_BRIDGE_H
#define LOBBY_BRIDGE_H

#include <QObject>
#include <QString>

// QML-callable bridge: tablet file ops -> Writerdeck-server over the unix socket.
class LobbyBridge : public QObject
{
    Q_OBJECT
public:
    void setRoot(QObject *root) { m_root = root; }

public slots:
    Q_INVOKABLE void requestNotesList();
    Q_INVOKABLE void createNote(const QString &name);
    Q_INVOKABLE void deleteNote(const QString &name);
    Q_INVOKABLE void renameNote(const QString &oldName, const QString &newName);
    Q_INVOKABLE void notifyOpen(const QString &name);
    Q_INVOKABLE void syncNow();
    Q_INVOKABLE void setKeyboardLayout(const QString &layout);

    // Called from the socket thread when the server pushes a notes list.
    void deliverNotesList(const QVariantList &items);

private:
    void sendReq(const QString &jsonLine);
    QObject *m_root = nullptr;
};

#endif
