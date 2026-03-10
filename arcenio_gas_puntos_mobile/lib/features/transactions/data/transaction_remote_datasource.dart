import 'package:dio/dio.dart';
import '../../../core/network/dio_client.dart';

class TransactionRemoteDatasource {
  final DioClient _dioClient;

  TransactionRemoteDatasource(this._dioClient);

  Future<Map<String, dynamic>> earnPoints(
    String clientId,
    double gallons,
  ) async {
    try {
      final response = await _dioClient.dio.post(
        '/transactions/earn',
        data: {'client_id': clientId, 'gallons': gallons},
      );
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(
          e.response?.data['error'] ?? 'Error al acumular puntos',
        );
      }
      throw Exception('Error de conexión con el servidor');
    }
  }

  Future<Map<String, dynamic>> redeemPoints(String clientId) async {
    try {
      final response = await _dioClient.dio.post(
        '/transactions/redeem',
        data: {'client_id': clientId},
      );
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(e.response?.data['error'] ?? 'Error al canjear puntos');
      }
      throw Exception('Error de conexión con el servidor');
    }
  }

  Future<Map<String, dynamic>> getBalance(String clientId) async {
    try {
      final response = await _dioClient.dio.get(
        '/transactions/balance/$clientId',
      );
      return response.data as Map<String, dynamic>;
    } on DioException catch (e) {
      if (e.response != null) {
        throw Exception(e.response?.data['error'] ?? 'Error al obtener saldo');
      }
      throw Exception('Error de conexión con el servidor');
    }
  }
}
