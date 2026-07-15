#include "lobby_bridge.h"

#include <QJsonDocument>
#include <QJsonObject>
#include <QMetaObject>
#include <QVariant>

// Defined in socket-inject.patch (main.cpp).
extern void rmkbdWriteLine(const std::string &line);

void LobbyBridge::sendReq(const QString &jsonLine)
{
    QByteArray ba = jsonLine.toUtf8();
    if (!ba.endsWith('\n'))
        ba.append('\n');
    rmkbdWriteLine(std::string(ba.constData(), static_cast<size_t>(ba.size())));
}

void LobbyBridge::requestNotesList()
{
    sendReq(QStringLiteral("{\"t\":\"req\",\"op\":\"noteslist\"}"));
}

void LobbyBridge::createNote(const QString &name)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("createnote");
    o[QStringLiteral("name")] = name;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::deleteNote(const QString &name)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("deletenote");
    o[QStringLiteral("name")] = name;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::renameNote(const QString &oldName, const QString &newName)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("renamenote");
    o[QStringLiteral("old")] = oldName;
    o[QStringLiteral("name")] = newName;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::notifyOpen(const QString &name)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("open");
    o[QStringLiteral("name")] = name;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::notifyReadOpen(const QString &name)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("openread");
    o[QStringLiteral("name")] = name;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::notifyLobbyInput(const QString &mode)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("lobbyinput");
    o[QStringLiteral("mode")] = mode;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::syncNow()
{
    sendReq(QStringLiteral("{\"t\":\"req\",\"op\":\"syncnow\"}"));
}

void LobbyBridge::setKeyboardLayout(const QString &layout)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("setkeyboardlayout");
    o[QStringLiteral("name")] = layout;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::setReadFont(const QString &font)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("setreadfont");
    o[QStringLiteral("name")] = font;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::setPinDigits(const QString &digits)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("setpindigits");
    o[QStringLiteral("name")] = digits;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::exitWriterdeck()
{
    sendReq(QStringLiteral("{\"t\":\"req\",\"op\":\"shutdown\"}"));
}

void LobbyBridge::setVaultPin(const QString &pin)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("setvaultpin");
    o[QStringLiteral("name")] = pin;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::changeVaultPin(const QString &oldPin, const QString &newPin)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("changevaultpin");
    o[QStringLiteral("old")] = oldPin;
    o[QStringLiteral("name")] = newPin;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::verifyVaultPin(const QString &pin, bool keepSession)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("verifyvaultpin");
    o[QStringLiteral("name")] = pin;
    o[QStringLiteral("old")] = keepSession ? QStringLiteral("session") : QStringLiteral("once");
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::encryptNote(const QString &name)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("encryptnote");
    o[QStringLiteral("name")] = name;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::decryptNote(const QString &name)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("decryptnote");
    o[QStringLiteral("name")] = name;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::createEncryptedNote(const QString &name)
{
    QString base = name.trimmed();
    if (base.isEmpty())
        return;
    if (!base.endsWith(QStringLiteral(".md.enc")))
        base += QStringLiteral(".md.enc");
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("req");
    o[QStringLiteral("op")] = QStringLiteral("createnote");
    o[QStringLiteral("name")] = base;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::publishState(int cursor, int selStart, int selEnd, int textLen, int mode, int isLobby,
                               const QString &vaultOverlay, const QString &currentFile,
                               const QString &text, int contentY)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("state");
    o[QStringLiteral("cursor")] = cursor;
    o[QStringLiteral("selStart")] = selStart;
    o[QStringLiteral("selEnd")] = selEnd;
    o[QStringLiteral("textLen")] = textLen;
    o[QStringLiteral("mode")] = mode;
    o[QStringLiteral("isLobby")] = isLobby;
    o[QStringLiteral("vaultOverlay")] = vaultOverlay;
    o[QStringLiteral("currentFile")] = currentFile;
    o[QStringLiteral("contentY")] = contentY;
    if (!text.isNull())
        o[QStringLiteral("text")] = text;
    sendReq(QString::fromUtf8(QJsonDocument(o).toJson(QJsonDocument::Compact)));
}

void LobbyBridge::deliverNotesList(const QVariantList &items)
{
    if (!m_root)
        return;
    QMetaObject::invokeMethod(m_root, "setNotesList",
        Qt::QueuedConnection,
        Q_ARG(QVariant, QVariant(items)));
}
