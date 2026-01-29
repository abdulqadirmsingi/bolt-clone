import 'package:flutter/material.dart';

import 'app/app.dart';
import 'app/di.dart';

void main() {
  // Start driver location streamer: only sends GPS to backend when driver is ONLINE/ON_TRIP
  AppContainer.driverLocationStreamer.start();
  runApp(const MobilityApp());
}
