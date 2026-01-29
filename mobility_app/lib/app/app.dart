import 'package:flutter/material.dart';

import '../features/map/presentation/map_page.dart';
import 'theme.dart';

class MobilityApp extends StatelessWidget {
  const MobilityApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Mobility',
      theme: MobilityTheme.light,
      darkTheme: MobilityTheme.dark,
      home: const MapPagePlaceholder(),
    );
  }
}

/// Placeholder until router is wired; replace with go_router home.
class MapPagePlaceholder extends StatelessWidget {
  const MapPagePlaceholder({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('Mobility', style: Theme.of(context).textTheme.headlineMedium),
            const SizedBox(height: 16),
            FilledButton(
              onPressed: () => Navigator.of(context).push(
                MaterialPageRoute(builder: (_) => const MapPage()),
              ),
              child: const Text('Open map'),
            ),
          ],
        ),
      ),
    );
  }
}
