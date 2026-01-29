import '../features/driver/data/driver_availability_api_noop.dart';
import '../features/driver/data/driver_location_sender.dart';
import '../features/driver/data/driver_location_streamer.dart';
import '../features/driver/data/driver_repository_impl.dart';
import '../features/driver/data/driver_stream_ds.dart';
import '../features/driver/domain/get_smoothed_driver_location.dart';
import '../features/driver/domain/driver_repository.dart';
import '../features/driver/presentation/driver_mode_controller.dart';
import '../features/map/domain/map_use_case.dart';
import '../features/map/presentation/map_controller.dart';

/// Simple dependency container: feature-first, data → domain → presentation.
///
/// Data layer uses core/location (smoothing); presentation never touches it.
class AppContainer {
  AppContainer._();

  static DriverStreamDataSource? _driverStreamDs;
  static DriverRepository? _driverRepository;
  static GetSmoothedDriverLocation? _getSmoothedDriverLocation;
  static ObserveDriverOnMap? _observeDriverOnMap;
  static MapController? _mapController;
  static DriverModeController? _driverModeController;
  static DriverLocationStreamer? _driverLocationStreamer;

  static DriverLocationSender get driverLocationSender =>
      DriverLocationSenderNoop();

  static DriverLocationStreamer get driverLocationStreamer {
    _driverLocationStreamer ??= DriverLocationStreamer(
      driverModeController: driverModeController,
      sender: driverLocationSender,
    );
    return _driverLocationStreamer!;
  }

  static void registerDriverStreamDataSource(DriverStreamDataSource ds) {
    _driverStreamDs = ds;
  }

  static DriverModeController get driverModeController {
    _driverModeController ??= DriverModeController(
      api: DriverAvailabilityApiNoop(),
    );
    return _driverModeController!;
  }

  static DriverRepository get driverRepository {
    _driverRepository ??= DriverRepositoryImpl(
      dataSource: _driverStreamDs ?? _StubDriverStreamDataSource(),
    );
    return _driverRepository!;
  }

  static GetSmoothedDriverLocation get getSmoothedDriverLocation {
    _getSmoothedDriverLocation ??= GetSmoothedDriverLocation(driverRepository);
    return _getSmoothedDriverLocation!;
  }

  static ObserveDriverOnMap get observeDriverOnMap {
    _observeDriverOnMap ??= ObserveDriverOnMap(driverRepository);
    return _observeDriverOnMap!;
  }

  static MapController get mapController {
    _mapController ??= MapController(observeDriverOnMap: observeDriverOnMap);
    return _mapController!;
  }

  static void reset() {
    _driverStreamDs = null;
    _driverRepository = null;
    _getSmoothedDriverLocation = null;
    _observeDriverOnMap = null;
    _mapController?.dispose();
    _mapController = null;
    _driverModeController?.dispose();
    _driverModeController = null;
    _driverLocationStreamer?.dispose();
    _driverLocationStreamer = null;
  }
}

/// Stub until gRPC is wired: no-op stream.
class _StubDriverStreamDataSource implements DriverStreamDataSource {
  @override
  Stream<RawDriverLocationUpdate> subscribeDriverLocation(String tripId, String driverId) =>
      const Stream.empty();
}
