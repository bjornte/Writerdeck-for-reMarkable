#include "rotation_watcher.h"
#include <QVariant>

void RotationWatcher::onRotationChanged()
{
    if (!m_root || !m_applyingFlag || !m_notify) return;
    if (*m_applyingFlag) return;
    m_notify(m_root->property("rotation").toInt());
}
