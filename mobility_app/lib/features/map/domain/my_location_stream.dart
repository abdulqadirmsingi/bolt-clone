import 'dart:async';

import 'package:geolocator/geolocator.dart';

/// Streams device location for map display. Does not send to backend.
/// Backend GPS stream (StreamDriverLocation) is only active when driver is ONLINE/ON_TRIP.
Stream<({double lat, double lng})> createMyLocationStream({
  double desiredAccuracy = LocationAccuracy.high,
  int distanceFilter = 5,
}) async* {
  final permission = await Geolocator.checkPermission();
  if (permission == LocationPermission.denied) {
    final requested = await Geolocator.requestPermission();
    if (requested != LocationPermission.whileInUse &&
        requested != LocationPermission.always) {
      return;
    }
  }
  if (permission == LocationPermission.deniedForever) return;

  final settings = LocationSettings(
    accuracy: desiredAccuracy,
    distanceFilter: distanceFilter,
  );
  await for (final pos in Geolocator.getPositionStream(locationSettings: settings)) {
    yield (lat: pos.latitude, lng: pos.longitude);
  }
}

/// Gets current position once (e.g. for initial map center).
Future<({double lat, double lng})?> getCurrentPositionOnce() async {
  final permission = await Geolocator.checkPermission();
  if (permission == LocationPermission.denied) {
    final requested = await Geolocator.requestPermission();
    if (requested != LocationPermission.whileInUse &&
        requested != LocationPermission.always) {
      return null;
    }
  }
  if (permission == LocationPermission.deniedForever) return null;
  try {
    final pos = await Geolocator.getCurrentPosition(
      desiredAccuracy: LocationAccuracy.high,
    );
    return (lat: pos.latitude, lng: pos.longitude);
  } catch (_) {
    return null;
  }
}
