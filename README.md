# 모바일 청첩장 서버

이 프로젝트는 [모바일 청첩장](https://github.com/juhonamnam/wedding-invitation) 웹 애플리케이션의 백엔드 서버입니다. 모바일 청첩장에 필요한 API 엔드포인트와 데이터베이스 관리 기능을 제공합니다. 모바일 청첩장에 필요한 간단한 기능만 구현하였으며, 트래픽이 많지 않은 환경이기에 SQLite를 사용합니다.

## 사전 요구사항

- Go 1.18

## 제공 기능

- 방명록 작성 및 조회 API
  - 관리자 비밀번호를 통한 방명록 강제 삭제 기능
- 참석 의사 전달 API
  - 참석 의사 제출 시 서버 로그(stdout)에 기록 → Railway 대시보드에서 확인 가능
  - 관리자 비밀번호를 통한 참석자 목록 조회 기능 (`GET /api/attendance?password=...`)

## 시작하기

1. 저장소 복제:

   ```bash
   git clone https://github.com/boseungjeong/wedding-invitation-server.git
   cd wedding-invitation-server
   ```

2. 의존성 설치:

   ```bash
   go mod download
   ```

3. 환경변수 설정:

   환경변수 샘플은 `.env.example` 파일에 저장되어 있습니다. 이 파일을 복사하여 `.env` 파일을 생성하고 각 환경변수를 수정합니다.

   ```bash
   cp .env.example .env
   ```

   - `ALLOW_ORIGIN`
     - 허용할 도메인
   - `ADMIN_PASSWORD`
     - 관리자 전용 비밀번호
     - 방명록 강제 삭제를 원하는 경우 해당 비밀번호로 삭제 가능
     - 참석자 목록 조회(`GET /api/attendance`) 시에도 동일 비밀번호 필요
   - `DB_PATH`
     - SQLite 데이터베이스 파일 경로 (기본값: `./sql.db`)
     - Railway Volume 사용 시 `/data/sql.db` 등 Volume 마운트 경로로 설정

## Railway Volume 설정 (데이터 영속화)

Railway의 기본 파일시스템은 배포/재시작 시 초기화됩니다. 방명록과 참석자 데이터를 유지하려면 Volume을 연결해야 합니다.

1. Railway 대시보드 → 서비스 선택 → **Volumes** 탭
2. **Add Volume** 클릭 → Mount Path를 `/data` 로 설정
3. **Variables** 탭에서 환경변수 추가: `DB_PATH=/data/sql.db`
4. 서비스를 재배포하면 이후 모든 데이터가 Volume에 저장되어 배포/재시작에도 유지됩니다.

## 관리자 기능

- 참석자 목록 조회:
  ```bash
  curl "https://<railway-domain>/api/attendance?password=<ADMIN_PASSWORD>"
  ```
  또는 헤더 방식:
  ```bash
  curl -H "X-Admin-Password: <ADMIN_PASSWORD>" "https://<railway-domain>/api/attendance"
  ```
- 참석 의사 제출 로그는 Railway 프로젝트 대시보드의 **Deployments → Logs** 탭에서 확인 가능합니다. `[ATTENDANCE]` 접두사로 필터링할 수 있습니다.

4. 서버 실행:
   ```bash
   go run app.go
   ```

   서버가 기본적으로 `http://localhost:8080`에서 실행됩니다.

## 배포하기

1. 프로젝트 빌드:
   ```bash
   go build
   ```

2. 빌드된 바이너리 실행:
   ```bash
   ./wedding-invitation-server
   ```
