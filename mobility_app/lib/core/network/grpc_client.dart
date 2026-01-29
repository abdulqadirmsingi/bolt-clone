import 'package:grpc/grpc.dart';

/// Shared gRPC channel: one persistent connection per app (no HTTP polling).
/// Configure with interceptors for JWT and metadata (OS, locale, app version).
class MobilityGrpcClient {
  late final ClientChannel _channel;
  late final Map<String, String> _metadata;

  MobilityGrpcClient({
    required String host,
    int port = 50051,
    bool isSecure = false,
    String? token,
    String? locale,
    String? appVersion,
  }) {
    _channel = ClientChannel(
      host,
      port: port,
      options: ChannelOptions(
        credentials: isSecure ? ChannelCredentials.secure() : ChannelCredentials.insecure(),
      ),
    );
    _metadata = {};
    if (token != null) _metadata['authorization'] = 'Bearer $token';
    if (locale != null) _metadata['x-locale'] = locale;
    if (appVersion != null) _metadata['x-app-version'] = appVersion;
  }

  ClientChannel get channel => _channel;
  Map<String, String> get metadata => Map.unmodifiable(_metadata);

  /// Call options with auth/metadata for each RPC.
  CallOptions get callOptions => CallOptions(metadata: _metadata);

  void close() => _channel.shutdown();
}
