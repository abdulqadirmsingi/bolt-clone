import '../domain/driver_entity.dart';

/// Presentation state: current driver position for map display.
///
/// Holds only smoothed position; no Kalman/dead reckoning in presentation.
class DriverState {
  const DriverState({
    this.position,
    this.error,
    this.isLoading = true,
  });

  final SmoothedDriverPosition? position;
  final String? error;
  final bool isLoading;

  DriverState copyWith({
    SmoothedDriverPosition? position,
    String? error,
    bool? isLoading,
  }) =>
      DriverState(
        position: position ?? this.position,
        error: error ?? this.error,
        isLoading: isLoading ?? this.isLoading,
      );
}
