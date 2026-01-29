import 'package:flutter/material.dart';

import 'di.dart';
import 'theme.dart';
import '../features/map/presentation/map_page.dart';

class MobilityApp extends StatelessWidget {
  const MobilityApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Mobility',
      theme: MobilityTheme.light,
      darkTheme: MobilityTheme.dark,
      debugShowCheckedModeBanner: false,
      home: MapPage(
        mapController: AppContainer.mapController,
        tripId: null,
        driverId: null,
      ),
    );
  }
}
