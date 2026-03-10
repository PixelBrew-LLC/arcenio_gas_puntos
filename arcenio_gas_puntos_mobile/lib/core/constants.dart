class AppConstants {
  AppConstants._();

  // Cambiar a la IP del servidor en la red local para producción
  static const String baseUrl =
      'http://192.168.1.7:3000'; // Android emulator -> localhost

  // Keys para secure storage
  static const String tokenKey = 'access_token';
  static const String loginDateKey = 'login_date';
  static const String userDataKey = 'user_data';
}
