#pragma once

#include <QObject>
#include <QString>

class QLocalServer;
class QLocalSocket;

// Minimal Writerdeck socket feed: NDJSON lines from the server into on-screen text.
class SocketFeed : public QObject {
    Q_OBJECT
    Q_PROPERTY(QString displayText READ displayText NOTIFY displayTextChanged)

public:
    explicit SocketFeed(QObject *parent = nullptr);
    ~SocketFeed() override;

    QString displayText() const { return m_text; }

    bool listen(const QString &path);

signals:
    void displayTextChanged();

private:
    void onNewConnection();
    void onReadyRead();
    void handleLine(const QByteArray &line);

    QLocalServer *m_server = nullptr;
    QLocalSocket *m_client = nullptr;
    QByteArray m_buffer;
    QString m_text;
};
