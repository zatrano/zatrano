<!DOCTYPE html>
<html lang="{{ str_replace('_', '-', app()->getLocale()) }}">
<head>
<title>ZATRANO | Giriş</title>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<!-- CSRF Token -->
<meta name="csrf-token" content="{{ csrf_token() }}">
<link rel="icon" type="image/png" href="{{ asset('assets') }}/login/images/icons/favicon.ico"/>
<link rel="stylesheet" type="text/css" href="{{ asset('assets') }}/login/fonts/iconic/css/material-design-iconic-font.min.css">
<link rel="stylesheet" type="text/css" href="{{ asset('assets') }}/login/css/util.css">
<link rel="stylesheet" type="text/css" href="{{ asset('assets') }}/login/css/main.css">
</head>
<body>
	<div class="limiter">
		<div class="container-login100" style="background-image: url('{{ asset('assets') }}/login/images/bg-01.jpg');">
			<div class="wrap-login100">
				<form class="login100-form validate-form" method="POST" action="{{ route('login') }}">
				@csrf
					<span class="login100-form-logo">
						<i class="zmdi zmdi-shield-check"></i>
					</span>
					<span class="login100-form-title p-b-34 p-t-27">
						ZATRANO
					</span>
					<div class="wrap-input100">
						<input id="email" type="email" class="input100 form-control @error('email') is-invalid @enderror" name="email" value="{{ old('email') }}" required autocomplete="email" placeholder="{{ __('Kullanıcı Adı') }}">
						@error('email')
							<span class="invalid-feedback" role="alert">
								<strong>{{ $message }}</strong>
							</span>
						@enderror
						<span class="focus-input100" data-placeholder="&#xf207;"></span>
					</div>
					<div class="wrap-input100">
						<input id="password" type="password" class="input100 form-control @error('password') is-invalid @enderror" name="password" required autocomplete="current-password" placeholder="{{ __('Password') }}">
						@error('password')
							<span class="invalid-feedback" role="alert">
								<strong>{{ $message }}</strong>
							</span>
						@enderror
						<span class="focus-input100" data-placeholder="&#xf191;"></span>
					</div>
					<div class="contact100-form-checkbox">
						<input class="input-checkbox100 form-check-input" type="checkbox" name="remember" id="remember" {{ old('remember') ? 'checked' : '' }} checked="">
						<label class="label-checkbox100 form-check-label" for="remember">
							Beni Hatırla
						</label>
					</div>
					<div class="container-login100-form-btn">
						<button type="submit" class="login100-form-btn">
							Giriş
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</body>
</html>
