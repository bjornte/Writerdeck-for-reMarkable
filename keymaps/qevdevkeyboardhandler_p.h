// Minimal QEvdevKeyboardMap definitions for kmap2qmap (from Qt 5.15 qtbase).
// Vendored so keymaps/generate.sh can build without QtInputSupport dev packages.

#ifndef QEVDEVKEYBOARDHANDLER_P_H
#define QEVDEVKEYBOARDHANDLER_P_H

#include <QtCore/qglobal.h>
#include <QtCore/qdatastream.h>
#include <QtCore/qnamespace.h>

namespace QEvdevKeyboardMap {

const quint32 FileMagic = 0x514d4150; // 'QMAP'

struct Mapping {
    quint16 keycode;
    quint16 unicode;
    quint32 qtcode;
    quint8 modifiers;
    quint8 flags;
    quint16 special;
};

enum Modifiers {
    ModPlain   = 0x00,
    ModShift   = 0x01,
    ModAltGr   = 0x02,
    ModControl = 0x04,
    ModAlt     = 0x08,
    ModShiftL  = 0x10,
    ModShiftR  = 0x20,
    ModCtrlL   = 0x40,
    ModCtrlR   = 0x80
};

enum Flags {
    IsDead     = 0x01,
    IsLetter   = 0x02,
    IsModifier = 0x04,
    IsSystem   = 0x08
};

enum SystemKeys {
    SystemConsoleFirst    = 0x0100,
    SystemConsoleLast     = 0x0107,
    SystemConsolePrevious = 0x0108,
    SystemConsoleNext     = 0x0109,
    SystemReboot          = 0x010a,
    SystemZap             = 0x010b
};

struct Composing {
    quint16 first;
    quint16 second;
    quint16 result;
};

inline QDataStream &operator>>(QDataStream &ds, Mapping &m)
{
    return ds >> m.keycode >> m.unicode >> m.qtcode >> m.modifiers >> m.flags >> m.special;
}

inline QDataStream &operator<<(QDataStream &ds, const Mapping &m)
{
    return ds << m.keycode << m.unicode << m.qtcode << m.modifiers << m.flags << m.special;
}

inline QDataStream &operator>>(QDataStream &ds, Composing &c)
{
    return ds >> c.first >> c.second >> c.result;
}

inline QDataStream &operator<<(QDataStream &ds, const Composing &c)
{
    return ds << c.first << c.second << c.result;
}

} // namespace QEvdevKeyboardMap

#endif
