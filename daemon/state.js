// state.js -- shared mutable application state.
// Both app.js and sync.js read and write these properties.
// Exporting a plain object rather than individual lets: property mutation
// crosses module boundaries freely (the binding -- the object reference --
// is constant; only properties change), sidestepping the read-only-binding
// rule on individually-exported lets with zero extra ceremony.
export var state = {
  syncOn: false,        // mirrors /api/settings syncOn
  syncRepo: '',         // mirrors /api/settings syncRepo
  tabletOpenNote: '',   // .md filename the tablet editor holds open; set by server
                        // openedit (phone /api/open or tablet doLoad); clears on exitedit
  typingMode: false     // false=Browse (list/read), true=Type (capture + echo)
};
