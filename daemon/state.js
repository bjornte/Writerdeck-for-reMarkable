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
  editorDiskHash: '',   // disk fingerprint at editor open — drift banner when disk changes (slice 8)
  typingMode: true,      // keyboard shell (capture + echo); default after auth
  remoteKeys: '',        // '' | 'read' | 'lobby' -- BT keyboard forward without full Type UI
  observeEnabled: false  // mirrors /api/settings observe -- phone Observe button
};
