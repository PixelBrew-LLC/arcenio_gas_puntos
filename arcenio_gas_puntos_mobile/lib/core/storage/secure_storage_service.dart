import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../constants.dart';

class SecureStorageService {
  final FlutterSecureStorage _storage = const FlutterSecureStorage();

  // --- Token ---
  Future<void> saveToken(String token) async {
    await _storage.write(key: AppConstants.tokenKey, value: token);
  }

  Future<String?> getToken() async {
    return await _storage.read(key: AppConstants.tokenKey);
  }

  Future<void> deleteToken() async {
    await _storage.delete(key: AppConstants.tokenKey);
  }

  // --- Login Date ---
  Future<void> saveLoginDate() async {
    final today = DateTime.now().toIso8601String().substring(
      0,
      10,
    ); // "2026-03-01"
    await _storage.write(key: AppConstants.loginDateKey, value: today);
  }

  Future<bool> isTokenValidForToday() async {
    final token = await getToken();
    if (token == null) return false;

    final loginDate = await _storage.read(key: AppConstants.loginDateKey);
    if (loginDate == null) return false;

    final today = DateTime.now().toIso8601String().substring(0, 10);
    return loginDate == today;
  }

  // --- User Data ---
  Future<void> saveUserData(Map<String, dynamic> userData) async {
    await _storage.write(
      key: AppConstants.userDataKey,
      value: jsonEncode(userData),
    );
  }

  Future<Map<String, dynamic>?> getUserData() async {
    final data = await _storage.read(key: AppConstants.userDataKey);
    if (data == null) return null;
    return jsonDecode(data) as Map<String, dynamic>;
  }

  // --- Clear All ---
  Future<void> clearAll() async {
    await _storage.deleteAll();
  }
}
