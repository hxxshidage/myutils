package ucry

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/hxxshidage/myutils/rand"
	"github.com/hxxshidage/myutils/time"
	"strings"
	"sync"
	"testing"
)

func TestHMAC_Encrypt(t *testing.T) {
	key := urand.RandStr(16)
	println(key)
	pwd := "super@Zhyadattr#521"
	println(HMAC.EncryptPlus(key, pwd, "md5"))
	println(utime.TsMill())
}

//func TestHMAC_Encrypt(t *testing.T) {
//	key := "dhvOYWmp7FsjGTqj"
//	println(key)
//	pwd := "123456"
//	// f1665cb04ec507179f486525715b492a
//	println(HMAC.EncryptPlus(key, pwd, "md5"))
//}

func TestRSA_EncryptPlus(t *testing.T) {
	// 1024
	//pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCXIoMwR734RAu4Xusz2sMyoovTSv4frlP16YfIVv4hyFnXtL4+94utUR4WVo+1WXpMuLSokD2w8CDaFR/C5zblSsxOmdKGkRcNRQrAacxxOkoSuJjLIHoehi292vBNdJTuyQKRBlzZnZRRC5MesPJLqDReZFspBBx6XJx/7EsIKwIDAQAB"
	//priKey := "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAJcigzBHvfhEC7he6zPawzKii9NK/h+uU/Xph8hW/iHIWde0vj73i61RHhZWj7VZeky4tKiQPbDwINoVH8LnNuVKzE6Z0oaRFw1FCsBpzHE6ShK4mMsgeh6GLb3a8E10lO7JApEGXNmdlFELkx6w8kuoNF5kWykEHHpcnH/sSwgrAgMBAAECgYB8Ngsn1O6WaiZPwwL/PR4MMXdFkm7EztuGUgYcaxK15RmhpJRu37hWG0LlDQNTAlT3VR51IwbwsontcksGPkzoY3X4AErdUPBKsRv8mj5iPTUkV2nx5P1HIZgJ5ztzdHF2iq2LNnLP9kF32BFA1SeidNP48LGbOKBwQiQG1NztsQJBAOGD1rxDW7VvI7d2H8fRCIt/J47gek+kW0S/AUFYgnefvKY7u2YF+jYUqfHVEr4b8V9Ml/vSQx4rn7f+a1XiOPcCQQCrkK04WrKI7dQ8BI8wRPF7eVK8OLuDfY7deUpfX4ZlG7L1Gd4w6ZyB0Da03KQ4Oiy2P41tdMZUE4vC6Tp7mrFtAkBjcUJog/9Vsyt5w+Hht8Bf2vMzorLsmZoRZ3SDPydQ3qROXTqWk34xCM2jsQRxxlqaXmkKhz5HeI8WkF7+YSIjAkBXsTl45pL+/mFO1B3EVFkB/b9WScE+snzFo2tqWE8/eur78N7rLV071QtBYs3ll4anGNDXM6rz6pGzbLbh2u35AkEAnuiSP65+ePGg4tGIN2bgfw6yNgGanDhsGPldp8d2/HdrQa5JGtLO2nfPpomJ6ga5BV5wK1jRnK3b6H+vQK7EKg=="
	pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC6d9F760l93kG/ii5NSXe4LlT4SPOlRb/5nfgiTgS3aD8VrphSzERzamk50y3iyTQb35E/UCcqjhcdfePaNrm+RiXWWAY21176BWi+WaO8yC7EHnRX3YZszuGGdS8YTzK0lIaUFv1Ulmo3Qe/A4SM2Cux4So6AZHY6pMegrrov/wIDAQAB"
	priKey := "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBALp30XvrSX3eQb+KLk1Jd7guVPhI86VFv/md+CJOBLdoPxWumFLMRHNqaTnTLeLJNBvfkT9QJyqOFx1949o2ub5GJdZYBjbXXvoFaL5Zo7zILsQedFfdhmzO4YZ1LxhPMrSUhpQW/VSWajdB78DhIzYK7HhKjoBkdjqkx6Cuui//AgMBAAECgYAzl7VheQ9jgRxl98Cp57r1PfSKofyv7OWaFkgyja2mZXSW9jD1L8l8uaOVuUWH6y28zJXuk6nSj2/72/owRVYMw4myiJZKEH8guadvbgxfXKVkMSnMB6q2PiYvYVLIRgZXNfs6puUWmZpypn/ROMqIAVMt1dyaOawJGbt06/bVsQJBANv/raFvDIQr2DYdLFqg854BV3j897w2xBLY6mkfTUmI+1mni01PLbgIKy5fjFUKiCEVYxWpX5D/xO40S/zsEmcCQQDY+3OMWiqx8Tn/Lk19Pet/zlezoHEX4lE0mDIAFmCwU2AHStRom/YCzngX4HST1iUAAK+gsXSZCgX8wwpF62apAkBT4MC3ixpvjjPqNedCSpl3xbUvGOvvY8YQFYoSOHAEOGj0vs9601RwFRU5og5kIhbS3k+f46rfItVXqLLMJMnDAkAs5swgjUAslpjD39cVRc0II7EucNM2eBUJ5zzhm7/ifT5wA+I0y6F7LzLW7hHxzezWNU4i2NYiHVrj1ZC4q2ARAkEAjwefHIQdelzVwSIgXVgLwbGTfl8rk6mAOBCeQYqZuEp7q5nKm4pnJ3vX59f/95ZEHm15fYGYqd5qHC0Hu+Pobw=="

	// 2048
	//pubKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0TZZ1P6lFM3eJ0sAuMczYFXzQiT5Rod8OjLi0m137QQ71kxeU4QT3fJ0OceNwmXlKHlqxQbNyj/H3gU6IWb0ui0zz6Nyqm+lf/tL8xsalSPwFmz0PxpqovOfatoyVxo6+zC2sMoll7bXF+MLsVsw/0/Kn/TYHsKzErfhlCs4L7WLDV44mygYzjXkKCuZZcpnlMzF5MB2y401LbHpZ+r31tTKdg69pg2xLI6UTUB/oHmENVvMPJ6SsGkvN1qDrza5GC6OjnrXDn+PqNIo/OA/p0ZZdRfwBjJtV+y6aIBfB0M5wUj3qNdr28gX3levy3s/5vW46LHirVWBiFP0fmn/qQIDAQAB"
	//priKey := "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDRNlnU/qUUzd4nSwC4xzNgVfNCJPlGh3w6MuLSbXftBDvWTF5ThBPd8nQ5x43CZeUoeWrFBs3KP8feBTohZvS6LTPPo3Kqb6V/+0vzGxqVI/AWbPQ/Gmqi859q2jJXGjr7MLawyiWXttcX4wuxWzD/T8qf9NgewrMSt+GUKzgvtYsNXjibKBjONeQoK5llymeUzMXkwHbLjTUtseln6vfW1Mp2Dr2mDbEsjpRNQH+geYQ1W8w8npKwaS83WoOvNrkYLo6OetcOf4+o0ij84D+nRll1F/AGMm1X7LpogF8HQznBSPeo12vbyBfeV6/Lez/m9bjoseKtVYGIU/R+af+pAgMBAAECggEAZrlsJGQgRHXM9bgjE5djx5KMTdb/ury9p4niy3XKo8snYlK/IfzBEIKCxPgRx/jmKxhq73EmzH7CsWYZo9r8oQme3f4gSEVnms6V/Tk6tS4fLbOzveRHpdk/VmTqwM/1U+8YVPf8u2Cgfm8SB7SB+2p/UEkVcQ0ihjdZgLoz7WAvuDuQ3NVeaHmoiG8CC6+b6ubvQJPKNkQRSS+VhAAS6LBr16P2A5QnGPvlkOYYA9dk0x16fK4dwnHH2fIiZNefwkcv3hM5A4eL79ymCpv/z2XiwPOiW5kCoc4lTX1FYCXA4bc4Cem5X800HjNKy4SRBzjK8c6YJGRQRDBT0v5HvQKBgQD214r2aL0/EgoeNJrV06A9FHlqEV38eXsQNmf5o950zVQvwFyRT/0m/KTKx0TFgEBwRelsAGxBsrzrcwRxa3OHKfoHaN4HwwoUaUeOXr76MuZCZTT32Y+lfQyWYUHVfI721CDzVIWjgIc/+Uuiij1G+RpUzJ+w2qamrplOBL0IfwKBgQDY+WitkbLUeEEeJz+Arv1/wZwQEcNu41oa/uc8LoM7hOGuWfRp1frQ3uT8VC04Cb7UEIX5Nc2xFLc3ksYIyjXUOTO9OfAopyQ+fyw8kt1aFH15+2XqG3Dr92nnvpL+cc1qnt+XjYNwRuYGJQ/Mga3COpLDlC3M9I/W0CVdgUGj1wKBgQC7ARyEBWGqEI6dx5it+f/hhktdcf7UMWxsIeuvktgLgSsQahk29XSCPtDR1xlgzwMCi9SFP+TD+RdrAN4S1ybU9ZY4Wtgq5TYTtJbDY0An6LvM/UOdqbNL2mrY2qG3jP/6O0cjUZtU/SB59PG6GCTIShwKtmMc5uluScRW+PwFHQKBgQCr0O0ErWC4gbXHI/tUcr4JQg7MZtSk5eJ/iCUicg63viVchJ1YhfsFFcysyBe6zXTQx3jf/Kwysx3XIyIw6bewo8+F4/B6sMixuNEV8pLYd2tZgiuVND+6jATQYAhU91dcPA0BHS3dZjdW3FhjcLlhGjMILzoJmAyjZdaA8g9BpwKBgDraf/xeK9bvIvLOqfb/5g+qtMp0xIYFqZccrH0J2O3T5dRp8GacHsSdOp79Bxzkf4DOrpFGEcOCz4RfPtSLX5IiOFldjJlYzKEZmyXqDInVW7QsTSAm7rxqjQ8e3kIit7BG6ll4Xqi+QvFKFGWOqET/fZzqZzN+wWzKhK1zBUkH"
	plain := "我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ"

	//plain := "super@@admin#"
	//plain := "1234567"
	//plain := "dage@gmail.com"
	val, err := RSA.EncryptPlus(pubKey, plain)
	println(val)
	if err != nil {
		println(err.Error())
	}

	val, err = RSA.DecryptPlus(priKey, val)
	println("=========")
	println(val)
	println(val == plain)
	if err != nil {
		println(err.Error())
	}
}

func BenchmarkRSA_Encrypt(b *testing.B) {
	// 1024
	//pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCrGh1sc5AKD1EQ8WdA1iWF4m7wXtO6WoS7Dtfd0Jm2ud+LKBQ+e7R6YIXnwfEKB/4Jm+jNtCi7/Zrx5gtEpUuVAyrEo5+qr5al5KibeJq3xyI/626IBsDMFX5o3WOoXceTF7+lgi6r+OuokqFJgpeh7YANXQ8Y8mn8ucw+Ly+LbQIDAQAB"
	//priKey := "MIICXQIBAAKBgQCrGh1sc5AKD1EQ8WdA1iWF4m7wXtO6WoS7Dtfd0Jm2ud+LKBQ+e7R6YIXnwfEKB/4Jm+jNtCi7/Zrx5gtEpUuVAyrEo5+qr5al5KibeJq3xyI/626IBsDMFX5o3WOoXceTF7+lgi6r+OuokqFJgpeh7YANXQ8Y8mn8ucw+Ly+LbQIDAQABAoGAGgoxbC3yP/WwyrlSk4WD1Gpvo9lqs7PO+4D4zWNP4YVMRitlWVUOVImYF3tmqbYprWCy/4tpn6KrECGImXvmkplXPxd4x3W+haZftx3VjTwh5fvT9yHp4swXxN+hLMItDdIOWS4U6wVJa77Dy7VfK303LZrPLqnxkf4oEywp5YECQQDZOz1WD7nOqOiyAlwDhfeLTmArN0f+gV6RLrxMp2XRqC2DN5nMq5O5BVVMK9LBgArNqYfxWYuMa3K2qliRDPPxAkEAyaNWq/fDvjpK9TgztqsHIiG+cUQpWI759zt5qHNA+QF4L43dtAVZzBR/uam1jnRuM6K0ZCSZo2ITiqapmk8bPQJAEd9d3IbOssIS4xJun5uWElAQeX3C3p2mOiuuMmBTcDx2AiXA8aXsMXzO18WDQYhXWzRniuPjJ1pvxbeeMdDvAQJBAMDhuZAJEzrOAlQurfFICyvQQZ+Rx0dKhbzFLOxBS96mVDSRLYn+MFbzKPcOa3lY0O4d7xd4l2td7zmLkePlVjUCQQCY8VuIfKc0+AWvPnktKXbx9bBdJZSDginZM5cu7pdxW0uB9KZoLqgbGLIvWrLyA6SBqo87Q1j1//wFgLP+A2Gn"

	pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC6d9F760l93kG/ii5NSXe4LlT4SPOlRb/5nfgiTgS3aD8VrphSzERzamk50y3iyTQb35E/UCcqjhcdfePaNrm+RiXWWAY21176BWi+WaO8yC7EHnRX3YZszuGGdS8YTzK0lIaUFv1Ulmo3Qe/A4SM2Cux4So6AZHY6pMegrrov/wIDAQAB"
	priKey := "MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBALp30XvrSX3eQb+KLk1Jd7guVPhI86VFv/md+CJOBLdoPxWumFLMRHNqaTnTLeLJNBvfkT9QJyqOFx1949o2ub5GJdZYBjbXXvoFaL5Zo7zILsQedFfdhmzO4YZ1LxhPMrSUhpQW/VSWajdB78DhIzYK7HhKjoBkdjqkx6Cuui//AgMBAAECgYAzl7VheQ9jgRxl98Cp57r1PfSKofyv7OWaFkgyja2mZXSW9jD1L8l8uaOVuUWH6y28zJXuk6nSj2/72/owRVYMw4myiJZKEH8guadvbgxfXKVkMSnMB6q2PiYvYVLIRgZXNfs6puUWmZpypn/ROMqIAVMt1dyaOawJGbt06/bVsQJBANv/raFvDIQr2DYdLFqg854BV3j897w2xBLY6mkfTUmI+1mni01PLbgIKy5fjFUKiCEVYxWpX5D/xO40S/zsEmcCQQDY+3OMWiqx8Tn/Lk19Pet/zlezoHEX4lE0mDIAFmCwU2AHStRom/YCzngX4HST1iUAAK+gsXSZCgX8wwpF62apAkBT4MC3ixpvjjPqNedCSpl3xbUvGOvvY8YQFYoSOHAEOGj0vs9601RwFRU5og5kIhbS3k+f46rfItVXqLLMJMnDAkAs5swgjUAslpjD39cVRc0II7EucNM2eBUJ5zzhm7/ifT5wA+I0y6F7LzLW7hHxzezWNU4i2NYiHVrj1ZC4q2ARAkEAjwefHIQdelzVwSIgXVgLwbGTfl8rk6mAOBCeQYqZuEp7q5nKm4pnJ3vX59f/95ZEHm15fYGYqd5qHC0Hu+Pobw=="

	// 2048
	//pubKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5Lrqv2/gTf6RJAAJMtyYWEWwLjW592rFCMIu6NX5/WR9/O1d3aPOHj6Zjt/cGaB3MBZp2J5nNj2HNGXjYB9mGI/JoLt4EQJo/2l2EBJPBNWkephDocBeenps2i7FGs1QvbPN0VYBkM9ArgndzCxSNZIbW+E8LI+l4Ri0e1uXS2XY3gv0mctdPkZuh90r1tSC+bqudBsLTIyhVpl/c4rNg5w1U0Ak3y4M1J8Jtl84hox5OwGV5Luk65V1UDeikU+jfFgmMvYKL3NI/CNbtwCPxfdrLOHi//Ls1zu5gdKF8io/xFN2hr76fb6HzmB8WQGjGrHWnkiYH29DlPjFCbFGkwIDAQAB"
	//priKey := "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDkuuq/b+BN/pEkAAky3JhYRbAuNbn3asUIwi7o1fn9ZH387V3do84ePpmO39wZoHcwFmnYnmc2PYc0ZeNgH2YYj8mgu3gRAmj/aXYQEk8E1aR6mEOhwF56emzaLsUazVC9s83RVgGQz0CuCd3MLFI1khtb4Twsj6XhGLR7W5dLZdjeC/SZy10+Rm6H3SvW1IL5uq50GwtMjKFWmX9zis2DnDVTQCTfLgzUnwm2XziGjHk7AZXku6TrlXVQN6KRT6N8WCYy9govc0j8I1u3AI/F92ss4eL/8uzXO7mB0oXyKj/EU3aGvvp9vofOYHxZAaMasdaeSJgfb0OU+MUJsUaTAgMBAAECggEAabU46GB7+Y+85DZgCfGJNsJ+Odz6pS3jAbk8lL7PWhwnXc0VpGkfyTqFHVK4Fd/jNYYmRMZwTNECu2SbQMFCHffV50K8qp/ChsfmmGbdvg4+han5F7gf8drCk9MppMlel02RwT1OW+5spgJJTyLsm3V6z2u4s59vuigAwUIDGgSID7+fwejS0gQUCMsrApEiG2nMvvAsjsNVD7XZEuBN0APuij+k0OAulzpkCgpdMwwfYJdxB90WYFOtFtnfbsz6bVTAn7dAIu9XgXwyQeUPHfAVC9cz29Ueuts6yM4a7Jxz9zvXPT5EW5Ms/C8u6aI7fP2jV1eLReTg24uLZJ+EAQKBgQD2xy+u9oDJY4O4xMD6cFSOL5UERab8YJdlDF15G5Hys8ohySTmzdUvY2yWAq1v0DVl0JxwX5DjwpF7uGuZvywnIgSfFsN+RYhojStKGlUK8+231JSeDpKRSEPbqIZp8p0h1DOwwGWk9OQY6FI7sA7ih3pCrrH4cXf3+VbQzlnoswKBgQDtRxML4SgNeXfftiJhss4VI9ezTbf7/27Hjm+cPt/Ulpn3Sm0RTPGMx4lQGP/A4/iGqeXga8I7kd97vEBrMtyuCD7apczZ27ji/vaMsSn0f9JdPBHEtuEDHWv8YLpLYAv2aw3vNxamUmzSCAg1XGVAyNQKw/NfduyOATl+cM1aoQKBgGNd/yhPX7o31PFQYHg3RQTfyfwnY77Z0fxBR14dqN32YRzLlo1NMltbiHy65UVRrD6sCmIBSSE81kHgF4uX+9piC0RX8S3mJ7AZr+WtxrKbWAwekB04tvHDDHflWwJMS9M0VAAG6KbMaRBSc9JO6R9z99nj6Aum/OyfvMJlZLSbAoGBAJUzjhx4NnFSojhAFRqODtxoL2iGRFznX8eIH1KGjsTk9mfzmuW4FmPJzORa8+dc8pfrGaum1voSXg82buN5lh6w/KUMgOW1Lms+m9YYSSN/hM4vyZSC0rbct1x5jmt7N8p5wsdbQpjPV7IybsbvFJRKNFuYn961r0YUKw0A7YBBAoGBAK4Qd3+VhWMsAdNJOdUbrJzk9t0Tz5AvlVjzLEiLkDdb0LqlY8WfvsROBRfRrgfihQJUrJe50VhMz9J0yEYNyQ0Lq5Kukt9fPbU4P8uf/8Dmk1c27aveVHGFTcslWRLEH3z6pn7QHG4sbjBLRqLD0vE96F79K8SkRQ5bAouzdUqE"

	//plain := "我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ"
	plain := "b8ddd173fdb07b9519917ff86b264b6e"

	for i := 0; i < b.N; i++ {
		val, err := RSA.EncryptPlus(pubKey, plain)
		//println(val)
		if err != nil {
			println(err.Error())
		}

		_, err = RSA.DecryptPlus(priKey, val)
		//println("=========")
		//println(val)
		//println(val == plain)
		if err != nil {
			println(err.Error())
		}
	}
}

func BenchmarkRSAHolder_EncryptPlus(b *testing.B) {
	// 1024
	//pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCrGh1sc5AKD1EQ8WdA1iWF4m7wXtO6WoS7Dtfd0Jm2ud+LKBQ+e7R6YIXnwfEKB/4Jm+jNtCi7/Zrx5gtEpUuVAyrEo5+qr5al5KibeJq3xyI/626IBsDMFX5o3WOoXceTF7+lgi6r+OuokqFJgpeh7YANXQ8Y8mn8ucw+Ly+LbQIDAQAB"
	//priKey := "MIICXQIBAAKBgQCrGh1sc5AKD1EQ8WdA1iWF4m7wXtO6WoS7Dtfd0Jm2ud+LKBQ+e7R6YIXnwfEKB/4Jm+jNtCi7/Zrx5gtEpUuVAyrEo5+qr5al5KibeJq3xyI/626IBsDMFX5o3WOoXceTF7+lgi6r+OuokqFJgpeh7YANXQ8Y8mn8ucw+Ly+LbQIDAQABAoGAGgoxbC3yP/WwyrlSk4WD1Gpvo9lqs7PO+4D4zWNP4YVMRitlWVUOVImYF3tmqbYprWCy/4tpn6KrECGImXvmkplXPxd4x3W+haZftx3VjTwh5fvT9yHp4swXxN+hLMItDdIOWS4U6wVJa77Dy7VfK303LZrPLqnxkf4oEywp5YECQQDZOz1WD7nOqOiyAlwDhfeLTmArN0f+gV6RLrxMp2XRqC2DN5nMq5O5BVVMK9LBgArNqYfxWYuMa3K2qliRDPPxAkEAyaNWq/fDvjpK9TgztqsHIiG+cUQpWI759zt5qHNA+QF4L43dtAVZzBR/uam1jnRuM6K0ZCSZo2ITiqapmk8bPQJAEd9d3IbOssIS4xJun5uWElAQeX3C3p2mOiuuMmBTcDx2AiXA8aXsMXzO18WDQYhXWzRniuPjJ1pvxbeeMdDvAQJBAMDhuZAJEzrOAlQurfFICyvQQZ+Rx0dKhbzFLOxBS96mVDSRLYn+MFbzKPcOa3lY0O4d7xd4l2td7zmLkePlVjUCQQCY8VuIfKc0+AWvPnktKXbx9bBdJZSDginZM5cu7pdxW0uB9KZoLqgbGLIvWrLyA6SBqo87Q1j1//wFgLP+A2Gn"

	pubKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCpm3UbxKYm+Nwp6jUlfUuBGhxsxtuM4StZ1+F6kWCESGirN/DOZ86FAoKXTeNpf/h4JxD3XnV3O5EpLHfbySSuusVCrT4NqKkwGzyM1zJClRTMyDCQpVcauHGCdwYZHnj4Xdtc12+vUcAhACewEld3Dx9QtlinRdMutkQobd8TGwIDAQAB"
	priKey := "MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBAKmbdRvEpib43CnqNSV9S4EaHGzG24zhK1nX4XqRYIRIaKs38M5nzoUCgpdN42l/+HgnEPdedXc7kSksd9vJJK66xUKtPg2oqTAbPIzXMkKVFMzIMJClVxq4cYJ3BhkeePhd21zXb69RwCEAJ7ASV3cPH1C2WKdF0y62RCht3xMbAgMBAAECgYBEseJMzOk+/6ysjV/ZP1ZFg/3fUOu7s7eLPBTnP9qHuYwrTQ0LjJ/o34tlHPbu1BYfFDOa/Xc2Q7oXoxsoud9Q/R4y5c7ZVR4uuaEHtUhwvK+VnY3k9k/dOAjtNAfWOwK/ZEa/Cz/tEaaop8eVtcmPaBGJgNg6Mily5yJadRsjgQJBANUprlpJrAaziOswyuvbNyRUD9XK2ceY846nG+Y2bPHKDKgsXyBAnWvz8rjIsm+NjT7veKajNkO2PLYDqKitAXcCQQDLsQc/CCQt2JsTgD+e9Il2ThBNQIBar+wkBZrn8poXw856Nw9xIxuanNzRed1pFXmdQWuO/2dXQcbJ4lrRfIR9AkAPw+BoPxChAkA7HMW1QeZHIox1RGZs86v3vfY7RYUzML1U5ss2SHEcHdOyxO0lgPOUVwO2V1XZFi5RS936c6krAkBvOJ21TDO9GN4sesXCfNImWB/MnuC3JAIz9R+NcUm0mkU/NJto8nubI/XrJ7i/LWu3c0ZQ0aLS4WazS7a9VlldAkAP64GYT0hbCkcAhN+b/jpqPUonvKibEFeNjO8jUnqlOXgWeYAVVW3vVIeJMq1AZrOo6eu1Z9TTYxpbLX+/AOED"

	// 2048
	//pubKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5Lrqv2/gTf6RJAAJMtyYWEWwLjW592rFCMIu6NX5/WR9/O1d3aPOHj6Zjt/cGaB3MBZp2J5nNj2HNGXjYB9mGI/JoLt4EQJo/2l2EBJPBNWkephDocBeenps2i7FGs1QvbPN0VYBkM9ArgndzCxSNZIbW+E8LI+l4Ri0e1uXS2XY3gv0mctdPkZuh90r1tSC+bqudBsLTIyhVpl/c4rNg5w1U0Ak3y4M1J8Jtl84hox5OwGV5Luk65V1UDeikU+jfFgmMvYKL3NI/CNbtwCPxfdrLOHi//Ls1zu5gdKF8io/xFN2hr76fb6HzmB8WQGjGrHWnkiYH29DlPjFCbFGkwIDAQAB"
	//priKey := "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDkuuq/b+BN/pEkAAky3JhYRbAuNbn3asUIwi7o1fn9ZH387V3do84ePpmO39wZoHcwFmnYnmc2PYc0ZeNgH2YYj8mgu3gRAmj/aXYQEk8E1aR6mEOhwF56emzaLsUazVC9s83RVgGQz0CuCd3MLFI1khtb4Twsj6XhGLR7W5dLZdjeC/SZy10+Rm6H3SvW1IL5uq50GwtMjKFWmX9zis2DnDVTQCTfLgzUnwm2XziGjHk7AZXku6TrlXVQN6KRT6N8WCYy9govc0j8I1u3AI/F92ss4eL/8uzXO7mB0oXyKj/EU3aGvvp9vofOYHxZAaMasdaeSJgfb0OU+MUJsUaTAgMBAAECggEAabU46GB7+Y+85DZgCfGJNsJ+Odz6pS3jAbk8lL7PWhwnXc0VpGkfyTqFHVK4Fd/jNYYmRMZwTNECu2SbQMFCHffV50K8qp/ChsfmmGbdvg4+han5F7gf8drCk9MppMlel02RwT1OW+5spgJJTyLsm3V6z2u4s59vuigAwUIDGgSID7+fwejS0gQUCMsrApEiG2nMvvAsjsNVD7XZEuBN0APuij+k0OAulzpkCgpdMwwfYJdxB90WYFOtFtnfbsz6bVTAn7dAIu9XgXwyQeUPHfAVC9cz29Ueuts6yM4a7Jxz9zvXPT5EW5Ms/C8u6aI7fP2jV1eLReTg24uLZJ+EAQKBgQD2xy+u9oDJY4O4xMD6cFSOL5UERab8YJdlDF15G5Hys8ohySTmzdUvY2yWAq1v0DVl0JxwX5DjwpF7uGuZvywnIgSfFsN+RYhojStKGlUK8+231JSeDpKRSEPbqIZp8p0h1DOwwGWk9OQY6FI7sA7ih3pCrrH4cXf3+VbQzlnoswKBgQDtRxML4SgNeXfftiJhss4VI9ezTbf7/27Hjm+cPt/Ulpn3Sm0RTPGMx4lQGP/A4/iGqeXga8I7kd97vEBrMtyuCD7apczZ27ji/vaMsSn0f9JdPBHEtuEDHWv8YLpLYAv2aw3vNxamUmzSCAg1XGVAyNQKw/NfduyOATl+cM1aoQKBgGNd/yhPX7o31PFQYHg3RQTfyfwnY77Z0fxBR14dqN32YRzLlo1NMltbiHy65UVRrD6sCmIBSSE81kHgF4uX+9piC0RX8S3mJ7AZr+WtxrKbWAwekB04tvHDDHflWwJMS9M0VAAG6KbMaRBSc9JO6R9z99nj6Aum/OyfvMJlZLSbAoGBAJUzjhx4NnFSojhAFRqODtxoL2iGRFznX8eIH1KGjsTk9mfzmuW4FmPJzORa8+dc8pfrGaum1voSXg82buN5lh6w/KUMgOW1Lms+m9YYSSN/hM4vyZSC0rbct1x5jmt7N8p5wsdbQpjPV7IybsbvFJRKNFuYn961r0YUKw0A7YBBAoGBAK4Qd3+VhWMsAdNJOdUbrJzk9t0Tz5AvlVjzLEiLkDdb0LqlY8WfvsROBRfRrgfihQJUrJe50VhMz9J0yEYNyQ0Lq5Kukt9fPbU4P8uf/8Dmk1c27aveVHGFTcslWRLEH3z6pn7QHG4sbjBLRqLD0vE96F79K8SkRQ5bAouzdUqE"

	//plain := "我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ"
	//plain := "woshi"
	plain := "woshi"

	rsaHolder, _ := NewRSAHolder(pubKey, priKey)

	for i := 0; i < b.N; i++ {
		val, err := rsaHolder.EncryptPlus(plain)
		//println(val)
		if err != nil {
			println(err.Error())
		}

		_, err = rsaHolder.DecryptPlus(val)
		//println("=========")
		//println(val)
		//println(val == plain)
		if err != nil {
			println(err.Error())
		}
	}
}

func TestAesEcbP7_EncryptPlus(t *testing.T) {
	//key := "yK4vSCuWZ2Bu4hOm"
	key := "sq6LwSwVUGPW12NBy3cm15Gd"
	//key := "y6fqzSX3kEEmsnRR1fWh8o3pDySwpspm"
	plain := "我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ我是hello world! jack 你好'@wtk' 呵呵...#$%$@ 有特殊字符码  大兄弟```~~~**EQ"
	val, err := AesEcbP7.EncryptPlus(key, plain)
	if err != nil {
		panic(err)
	}
	println(val)
	val, err = AesEcbP7.DecryptPlus(key, val)
	if err != nil {
		panic(err)
	}
	println(val)

	println(val == plain)
}

func TestAesCbcP7_Encrypt(t *testing.T) {
	//iv := "yK4vSCuWZ2Bu4hOm"
	//key := "yK4vSCuWZ2Bu4hOm"
	//key := "sq6LwSwVUGPW12NBy3cm15Gd"
	//key := "y6fqzSX3kEEmsnRR1fWh8o3pDySwpspm"

	key := "b8ddd173fdb07b9519917ff86b264b6e"
	iv := key[16:]

	plain := `
{
	"ref_url": "zhyad_bridge=true&aid=be81c346b3a7411ea17e205fcfdffb2b&fid=67f48634d3a8f93704a600d5&lpcts=1743660370994&fpid=67edef14ddcce491f48149ae&pid=67eb9e3f4505018367bbb2e4&adtype=3&lpid=67ee25528e91513eec55ce34&ptype=1",
    "ref_click_ts": 123,
    "app_install_ts": 23,
    "install_exp_launched": false
}
`
	val, err := AesCbcP7.EncryptPlus(key, iv, plain)
	if err != nil {
		panic(err)
	}
	println(val)
	val, err = AesCbcP7.DecryptPlus(key, iv, val)
	if err != nil {
		panic(err)
	}
	println(val)

	println(val == plain)
}

func TestHAsh_Encrypt(t *testing.T) {
	//println(HASH.EncryptPlus("1", "md5"))
	//println(HASH.EncryptPlus("1", "sha1"))
	//println(HASH.EncryptPlus("1", "sha256"))
	//println(HASH.EncryptPlus("1", "sha512"))
	//println(HMAC.EncryptPlus("1", "1", "sha256"))

	// 205fbf96a5b5667d08039621ae3615984e7fc212
	// b90b32dec92b2cb7a8f65713e9e14cd60ab865b8
	content := `{"attr_id":"e3a2d0abad6b403b9c4e78e7e96fc9c0","ua":"mozilla/5.0 (windows nt 10.0; win64; x64) applewebkit/537.36 (khtml, like gecko) chrome/134.0.0.0 safari/537.36","clipboard":"","ad_flow_id":"67f48634d3a8f93704a600d5"}`
	println(HMAC.EncryptPlus("67f48634d3a8f93704a600d5", content, "sha1"))

}

func TestSyncMap(t *testing.T) {
	var m sync.Map

	val, ok := m.LoadOrStore("1", 111)
	println(val.(int), ok)
	val, ok = m.LoadOrStore("1", 111)
	println(val.(int), ok)
}

func TestGenKeys(t *testing.T) {
	priKey, err2 := rsa.GenerateKey(rand.Reader, 1024)
	if err2 != nil {
		panic(err2)
	}

	//derStream, err := x509.MarshalPKCS8PrivateKey(priKey)
	derStream := x509.MarshalPKCS1PrivateKey(priKey)

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	prvKey := pem.EncodeToMemory(block)
	puKey := &priKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(puKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubKey := pem.EncodeToMemory(block)

	publicKey := string(pubKey)
	privateKey := string(prvKey)

	println("pubKey: " + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(publicKey, "\n", ""), "-----BEGIN PUBLIC KEY-----", ""), "-----END PUBLIC KEY-----", ""))
	println("priKey: " + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(privateKey, "\n", ""), "-----BEGIN PRIVATE KEY-----", ""), "-----END PRIVATE KEY-----", ""))
	//println(strings.ReplaceAll(publicKey, "\n", ""))
}

func TestGenAk(t *testing.T) {
	println(urand.RandStr(32))
}
