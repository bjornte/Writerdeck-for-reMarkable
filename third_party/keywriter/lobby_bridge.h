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
    Q_INVOKABLE void notifyReadOpen(const QString &name);
    Q_INVOKABLE void notifyLobbyInput(const QString &mode);
    Q_INVOKABLE void syncNow();
    Q_INVOKABLE void setKeyboardLayout(const QString &layout);
    Q_INVOKABLE void setReadFont(const QString &font);
    Q_INVOKABLE void setPinDigits(const QString &digits);
    Q_INVOKABLE void exitWriterdeck();
    Q_INVOKABLE void setVaultPin(const QString &pin);
    Q_INVOKABLE void changeVaultPin(const QString &oldPin, const QString &newPin);
    Q_INVOKABLE void verifyVaultPin(const QString &pin, bool keepSession);
    Q_INVOKABLE void encryptNote(const QString &name);
    Q_INVOKABLE void decryptNote(const QString &name);
    Q_INVOKABLE void createEncryptedNote(const QString &name);
    Q_INVOKABLE void publishState(int cursor, int selStart, int selEnd, int textLen, int mode, int isLobby,
                                  const QString &vaultOverlay, const QString &currentFile,
                                  const QString &text = QString(), int contentY = 0);

    // Called from the socket thread when the server pushes a notes list.
    void deliverNotesList(const QVariantList &items);

private:
    void sendReq(const QString &jsonLine);
    QObject *m_root = nullptr;
};

#endif
