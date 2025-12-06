PaymentIntent 後続ユースケース TODO
- CapturePaymentIntentUseCase: オーソリ済み意図を売上計上し、succeeded / requires_payment_method / canceled へ遷移
- FailPaymentIntentUseCase: 決済失敗時のイベント記録、requires_payment_method へ戻すか canceled に遷移
- CancelPaymentIntentUseCase: ユーザー／システム起因のキャンセルで canceled に遷移
- RefundPaymentIntentUseCase (必要なら): 部分／全額返金を記録し残高管理
- ApplyPaymentIntentEventUseCase: Webhook やキュー経由の外部イベントを検証し、Intent に適用
- Converter 更新: 新しいステータスやフィールドを反映するため PaymentIntent view 変換を拡張
