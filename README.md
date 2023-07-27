# account example
- golang 숙련 및 테스트를 위한 Repository.
- Account에 대한 crud 동작 구현
- Account정보 저장할 DB와의 통신 구현

### Account
- email
- password
    - use PBKDF2
- name
- birth
- gender

# 기타 조건
- web server를 통한 UI를 제공해야 함
- DB는 MongoDB를 사용
- Docker를 통해 빌드 - 배포 할 수 있어야 함
- 테스트 코드 작성