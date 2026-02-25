package fnMailer

import "fmt"

func GetInquiryAnswerTemplate(title, answer string) string {
	return fmt.Sprintf(`
		<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #eee;">
			<h2 style="color: #333;">1:1 문의하신 내용에 대한 답변입니다.</h2>
			<p style="color: #555; font-size: 14px;">안녕하세요. 고객님께서 문의하신 <b>[%s]</b> 에 대한 답변이 등록되었습니다.</p>
			
			<div style="background-color: #f9f9f9; padding: 15px; margin: 20px 0; border-radius: 5px; white-space: pre-wrap;">
%s
			</div>
			
			<p style="color: #999; font-size: 12px;">본 메일은 발신 전용이므로 회신되지 않습니다.</p>
		</div>
	`, title, answer)
}

func GetVerificationTemplate(code string) string {
	return fmt.Sprintf(`<h2>인증 코드: <strong>%s</strong></h2>`, code)
}
