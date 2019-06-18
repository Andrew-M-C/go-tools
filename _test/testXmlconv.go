package _test

import (
	"github.com/Andrew-M-C/go-tools/xmlconv"
	"github.com/Andrew-M-C/go-tools/log"
)

const _TEST_STR = `
<xml>
  <appid><![CDATA[wx2421b1c4370ec43b]]></appid>

  <attach><![CDATA[支付测试]]></attach>
  <bank_type><![CDATA[CFT]]></bank_type>
  <fee_type><![CDATA[CNY]]></fee_type>
  <is_subscribe><![CDATA[Y]]></is_subscribe>
  <mch_id><![CDATA[10000100]]></mch_id>
  <nonce_str><![CDATA[5d2b6c2a8db53831f7eda20af46e531c]]></nonce_str>
  <openid><![CDATA[oUpF8uMEb4qRXf22hE3X68TekukE]]></openid>
  <out_trade_no><![CDATA[<1409811653]]></out_trade_no>
  <result_code><![CDATA[SUCCESS]]></result_code>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <sign><![CDATA[B552ED6B279343CB493C5DD0D78AB241]]></sign>
  <sub_mch_id><![CDATA[10000100]]></sub_mch_id>
  <time_end><![CDATA[20140903131540]]></time_end>
  <total_fee>1</total_fee>
<coupon_fee_0><![CDATA[10]]></coupon_fee_0>
<coupon_count><![CDATA[1]]></coupon_count>
<coupon_type><![CDATA[CASH]]></coupon_type>
<coupon_id><![CDATA[10000]]></coupon_id>
  <trade_type><![CDATA[JSAPI]]></trade_type>
  <transaction_id><![CDATA[1004400740201409030005092168]]></transaction_id>
  <sub>
    <hello>hello</hello>
  </sub>
</xml>
`

func TestXmlconv() {
	log.Info("====== Now test xmlconv")
	defer log.Info("done\n")

	obj, err := xmlconv.NewFromString(_TEST_STR)
	if err != nil {
		log.Error("Failed to convert: %v", err)
		return
	}

	s, err := obj.Marshal()
	if err != nil {
		log.Error("Failed to marshal: %v", err)
		return
	}
	log.Info("re-marshal result: %s", s)

	obj, err = xmlconv.NewFromString(s)
	if err != nil {
		log.Error("re-unmarshal error: %v", err)
		return
	}
	log.Info("re-marshal OK")

	s, err = obj.Marshal(xmlconv.Option{Indent: "\t"})
	if err != nil {
		log.Error("Failed to marshal: %v", err)
		return
	}
	log.Info("re-marshal result: %s", s)

	return
}


func TestXmlconvAccess() {
	root := xmlconv.NewItem("root")
	root.SetEmptyChild("child_1", "child_2")
	root.SetChildString("hello, xml!", "child_A", "child_2", "child_3")
	root.SetEmptyChild("child_A", "child_2", "child_3")

	c := xmlconv.NewItem("")
	c.SetString("Hello, Earth!")
	c.SetAttr("attr_1", "\"1\"")
	c.SetAttr("attr_2", "2")
	root.SetChild(c, "child_A", "child_2", "child_3")
	s, _ := root.Marshal(xmlconv.Option{Indent:"  "})
	log.Info("marshal result: \n%s", s)

	// re-unmarshal
	root, _ = xmlconv.NewFromString(s)
	s, _ = root.Marshal(xmlconv.Option{Indent:""})
	log.Info("re-marshal result: \n%s", s)
	return
}
