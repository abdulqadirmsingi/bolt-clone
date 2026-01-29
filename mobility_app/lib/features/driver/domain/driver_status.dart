/// Driver presence/availability state.
/// State machine: OFFLINE → ONLINE → ON_TRIP → ONLINE → OFFLINE.
enum DriverStatus {
  offline,
  online,
  onTrip,
}

extension DriverStatusX on DriverStatus {
  bool get isAvailableForDiscovery => this == DriverStatus.online || this == DriverStatus.onTrip;
  String get label {
    switch (this) {
      case DriverStatus.offline:
        return 'Offline';
      case DriverStatus.online:
        return 'Online';
      case DriverStatus.onTrip:
        return 'On trip';
    }
  }
}
