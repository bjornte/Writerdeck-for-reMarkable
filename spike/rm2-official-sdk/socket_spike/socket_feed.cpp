#include "socket_feed.h"

#include <QJsonDocument>
#include <QJsonObject>
#include <QLocalServer>
#include <QLocalSocket>
#include <QFile>

SocketFeed::SocketFeed(QObject *parent) : QObject(parent) {}

SocketFeed::~SocketFeed() {
    if (m_server)
        m_server->close();
}

bool SocketFeed::listen(const QString &path) {
    QFile::remove(path);
    m_server = new QLocalServer(this);
    connect(m_server, &QLocalServer::newConnection, this, &SocketFeed::onNewConnection);
    if (!m_server->listen(path)) {
        return false;
    }
    return true;
}

void SocketFeed::onNewConnection() {
    if (m_client) {
        m_client->disconnect(this);
        m_client->deleteLater();
    }
    m_client = m_server->nextPendingConnection();
    if (!m_client)
        return;
    m_buffer.clear();
    connect(m_client, &QLocalSocket::readyRead, this, &SocketFeed::onReadyRead);
}

void SocketFeed::onReadyRead() {
    if (!m_client)
        return;
    m_buffer.append(m_client->readAll());
    int idx = 0;
    while ((idx = m_buffer.indexOf('\n')) >= 0) {
        const QByteArray line = m_buffer.left(idx);
        m_buffer.remove(0, idx + 1);
        if (!line.isEmpty())
            handleLine(line);
    }
}

void SocketFeed::handleLine(const QByteArray &line) {
    const QJsonDocument doc = QJsonDocument::fromJson(line);
    if (!doc.isObject())
        return;
    const QJsonObject obj = doc.object();
    const QString type = obj.value(QStringLiteral("t")).toString();
    if (type == QStringLiteral("text")) {
        const int cp = obj.value(QStringLiteral("cp")).toInt();
        if (cp > 0)
            m_text += QChar(cp);
        emit displayTextChanged();
        return;
    }
    if (type == QStringLiteral("key")) {
        const QString key = obj.value(QStringLiteral("k")).toString();
        if (key == QStringLiteral("Return"))
            m_text += QChar(u'\n');
        else if (key == QStringLiteral("Backspace") && !m_text.isEmpty())
            m_text.chop(1);
        emit displayTextChanged();
    }
}
