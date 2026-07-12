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

void LobbyBridge::publishState(int cursor, int selStart, int selEnd, int textLen, int mode)
{
    QJsonObject o;
    o[QStringLiteral("t")] = QStringLiteral("state");
    o[QStringLiteral("cursor")] = cursor;
    o[QStringLiteral("selStart")] = selStart;
    o[QStringLiteral("selEnd")] = selEnd;
    o[QStringLiteral("textLen")] = textLen;
    o[QStringLiteral("mode")] = mode;
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
