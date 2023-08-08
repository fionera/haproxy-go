// Code generated by "stringer -type ControlMessageType"; DO NOT EDIT.

package peers

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ControlMessageSyncRequest-0]
	_ = x[ControlMessageSyncFinished-1]
	_ = x[ControlMessageSyncPartial-2]
	_ = x[ControlMessageSyncConfirmed-3]
	_ = x[ControlMessageHeartbeat-4]
}

const _ControlMessageType_name = "ControlMessageSyncRequestControlMessageSyncFinishedControlMessageSyncPartialControlMessageSyncConfirmedControlMessageHeartbeat"

var _ControlMessageType_index = [...]uint8{0, 25, 51, 76, 103, 126}

func (i ControlMessageType) String() string {
	if i >= ControlMessageType(len(_ControlMessageType_index)-1) {
		return "ControlMessageType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ControlMessageType_name[_ControlMessageType_index[i]:_ControlMessageType_index[i+1]]
}
