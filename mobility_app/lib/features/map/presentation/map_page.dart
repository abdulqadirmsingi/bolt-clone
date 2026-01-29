import 'package:flutter/material.dart';

/// Map-first design: live map with polylines and animated car markers.
/// Uses Google Maps or Mapbox; bottom sheet for trip status and ETA.
class MapPage extends StatelessWidget {
  const MapPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          // Map (Google Maps / Mapbox) with polylines and markers
          const Center(child: Text('Map view (integrate Google Maps / Mapbox)')),
          // Bottom sheet: trip status, driver ETA, multi-stop progress
          Align(
            alignment: Alignment.bottomCenter,
            child: _TripBottomSheet(
              etaMinutes: 5,
              driverName: 'Driver',
              stopLabel: 'Pickup',
            ),
          ),
        ],
      ),
    );
  }
}

class _TripBottomSheet extends StatelessWidget {
  final int etaMinutes;
  final String driverName;
  final String stopLabel;

  const _TripBottomSheet({
    required this.etaMinutes,
    required this.driverName,
    required this.stopLabel,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.1),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(stopLabel, style: Theme.of(context).textTheme.titleMedium),
            const SizedBox(height: 8),
            Text('$driverName • ETA $etaMinutes min', style: Theme.of(context).textTheme.bodyLarge),
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: FilledButton(
                onPressed: () {},
                child: const Text('View trip'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
