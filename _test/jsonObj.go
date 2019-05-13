package _test

import (
	"github.com/Andrew-M-C/go-tools/jsonconv"
	"github.com/Andrew-M-C/go-tools/log"
	"strconv"
)

var strStandard = `{
	"a-string": "这是一个string",
	"an-int": 12345678,
	"a-float": 12345.12345678,
	"a-true": true,
	"a-false": false,
	"a-null": null,
	"an-object": {
		"sub-string": "string in an object",
		"sub-object": {
			"another-sub-string": "\"string\" in an object in an object",
			"another-sub-array": [1, "string in sub
array", true, null],
			"complex":"\u4e2d\t\u6587"
		}
	},
	"an-array": [
		{
			"sub-string": "string in an object in an array",
			"sub-sub-array": [
				{
					"sub-sub-string": "string in an object in an array in an object in an string"
				}
			]
		},
		56789,
		false,
		null
	]
}`

func testKeyInObject(obj *jsonconv.JsonValue, key interface{}, keys... interface{}) {
	log.Info("Test %v %v", key, keys)
	child, err := obj.Get(key, keys...)
	if err != nil {
		log.Error("Failed to get child: %s", err.Error())
	} else {
		if child.IsString() {
			log.Info("Get child, type: %s, value '%s'", child.TypeString(), child.String())
		} else {
			log.Info("Get child, type: %s", child.TypeString())
		}
	}
}

func TestJsonValue() {
	log.Info("======== Start testing JsonValue")
	obj, err := jsonconv.NewFromString(strStandard)
	if err != nil {
		log.Error("Failed to parse json: %s", err.Error())
	} else {
		testKeyInObject(obj, "an-object", "sub-object", "complex")
		// testKeyInObject(obj, "an-object", "sub-object", "illegal")
		testKeyInObject(obj, "an-array", 0, "sub-string")		// ERROR OCCURRED !!!
		testKeyInObject(obj, "an-object", "sub-object", "another-sub-array", 1)
	}

	json_str := ""
	json_str, _ = obj.Marshal()
	log.Info("re-package json: %s", json_str)
	json_str, _ = obj.Marshal(jsonconv.Option{EnsureAscii: true})
	log.Info("re-package json: %s", json_str)
	json_str, _ = obj.Marshal(jsonconv.Option{FloatDigits: 2})
	log.Info("re-package json: %s", json_str)

	// test foreach
	log.Info("Now test object foreach")
	obj.ObjectForeach(func (key string, value *jsonconv.JsonValue) error {
		log.Info("Key - %s, type %s", key, value.TypeString())
		return nil
	})

	log.Info("Now test array foreach")
	arr, err := obj.Get("an-array")
	log.Info("array type: %s", arr.TypeString())
	if err != nil {
		log.Error("error: %s", err.Error())
	} else {
		arr.ArrayForeach(func (index int, value *jsonconv.JsonValue) error {
			log.Info("[%d] type %s", index, value.TypeString())
			return nil
		})
	}

	// test modification
	err = obj.Set(jsonconv.NewString("THIS IS A FULL NEW STRING"), "an-array", 0, "sub-string")
	if err != nil {
		log.Error("Failed to set: %s", err.Error())
	}
	err = obj.Set(jsonconv.NewString("THIS IS ANOTHER FULL NEW STRING"), "an-array", 0, "sub-sub-array", 0, "sub-sub-string")
	if err != nil {
		log.Error("Failed to set: %s", err.Error())
	}
	json_str, _ = obj.Marshal()
	log.Info("new json after modification: %s", json_str)

	// test simple json string
	{
		a_str_value, err := jsonconv.NewFromString(`"hello"`)
		if err != nil {
			log.Error("unmarshal string failed: %s", err.Error())
			return
		}
		if a_str_value.Type() != jsonconv.String {
			log.Error("json type error: %s", a_str_value.TypeString())
			return
		}
		raw_string, _ := a_str_value.Marshal()
		log.Info("test basic string OK, value: %s", raw_string)
	}

	// test simple json bool
	{
		a_bool_value, err := jsonconv.NewFromString("true")
		if err != nil {
			log.Error("unmarshal string failed: %s", err.Error())
			return
		}
		if a_bool_value.Type() != jsonconv.Boolean {
			log.Error("json type error: %s", a_bool_value.TypeString())
			return
		}
		raw_string, _ := a_bool_value.Marshal()
		log.Info("test basic string OK, value: %s", raw_string)
	}

	// test simple null
	{
		a_null_value, err := jsonconv.NewFromString("null")
		if err != nil {
			log.Error("unmarshal string failed: %s", err.Error())
			return
		}
		if a_null_value.Type() != jsonconv.Null {
			log.Error("json type error: %s", a_null_value.TypeString())
			return
		}
		raw_string, _ := a_null_value.Marshal()
		log.Info("test basic string OK, value: %s", raw_string)
	}

	// test simple number
	{
		an_int_value, err := jsonconv.NewFromString("-1000.0001")
		if err != nil {
			log.Error("unmarshal string failed: %s", err.Error())
			return
		}
		if an_int_value.Type() != jsonconv.Number {
			log.Error("json type error: %s", an_int_value.TypeString())
			return
		}
		raw_string, _ := an_int_value.Marshal()
		log.Info("test basic string OK, value: %s", raw_string)
	}

	// test uint64
	{
		uint_val := uint64(0xFFFFFFFFFFFFFFFF)
		a_uint_value, err := jsonconv.NewFromString(strconv.FormatUint(uint_val, 10))
		if err != nil {
			log.Error("unmarshal string failed: %s", err.Error())
			return
		}
		if a_uint_value.Type() != jsonconv.Number {
			log.Error("json type error: %s", a_uint_value.TypeString())
			return
		}
		raw_string, _ := a_uint_value.Marshal()
		if raw_string == strconv.FormatUint(uint_val, 10) {
			log.Info("test basic string OK, value: %s", raw_string)
		} else {
			log.Error("could not recognize large uint64")
		}
	}

	// return
	return
}


func TestAwsomeEscapingJson() {
	raw := `{"code":"0","msg":"操作成功","data":"{\"result\":\"{\\\"pageNo\\\":1,\\\"pageSize\\\":20,\\\"maxPageSize\\\":100,\\\"totalCount\\\":1,\\\"resultList\\\":[{\\\"orderId\\\":624573044000041,\\\"srcOrderId\\\":\\\"624573044000041\\\",\\\"srcPlatId\\\":4,\\\"srcOrderType\\\":0,\\\"srcInnerType\\\":0,\\\"srcInnerOrderId\\\":0,\\\"orderType\\\":10000,\\\"orderStatus\\\":90000,\\\"orderStatusTime\\\":\\\"2016-10-11 10:18:11\\\",\\\"orderStartTime\\\":\\\"2016-10-11 09:50:43\\\",\\\"orderPurchaseTime\\\":\\\"2016-10-11 09:50:57\\\",\\\"orderAgingType\\\":12,\\\"orderPreStartDeliveryTime\\\":\\\"2016-10-11 11:50:00\\\",\\\"orderPreEndDeliveryTime\\\":\\\"2016-10-11 11:50:00\\\",\\\"orderIsClosed\\\":true,\\\"orderCloseTime\\\":\\\"2016-10-11 10:18:11\\\",\\\"orgCode\\\":\\\"74554\\\",\\\"buyerPinType\\\":0,\\\"buyerPin\\\":\\\"JD_350u24a96522f\\\",\\\"buyerFullName\\\":\\\"王斌\\\",\\\"buyerFullAddress\\\":\\\"武汉市江汉区金雅公寓2栋1单元902室\\\",\\\"buyerMobile\\\":\\\"13720339384\\\",\\\"buyerProvince\\\":\\\"0\\\",\\\"buyerCity\\\":\\\"1381\\\",\\\"buyerCountry\\\":\\\"3582\\\",\\\"produceStationNo\\\":\\\"10055023\\\",\\\"produceStationName\\\":\\\"可多直营-汉兴小区店\\\",\\\"produceStationNoIsv\\\":\\\"0014\\\",\\\"deliveryStationNo\\\":\\\"10055023\\\",\\\"deliveryStationName\\\":\\\"可多直营-汉兴小区店\\\",\\\"deliveryStationNoIsv\\\":\\\"0014\\\",\\\"deliveryType\\\":1,\\\"deliveryCarrierNo\\\":\\\"9966\\\",\\\"deliveryCarrierName\\\":\\\"达达专送\\\",\\\"deliveryBillNo\\\":\\\"624573044000041\\\",\\\"deliveryPackageWeight\\\":1.7179999649524689,\\\"deliveryPackageVolume\\\":0,\\\"deliveryManName\\\":\\\"涂志学\\\",\\\"deliveryManPhone\\\":\\\"15926351698\\\",\\\"deliveryConfirmTime\\\":\\\"2016-10-11 10:18:10\\\",\\\"orderPayType\\\":4,\\\"orderTotalMoney\\\":7300,\\\"orderDiscountMoney\\\":500,\\\"orderFreightMoney\\\":200,\\\"orderBuyerPayableMoney\\\":7000,\\\"orderVenderChargeMoney\\\":0,\\\"packagingMoney\\\":0,\\\"orderBalanceUsed\\\":0,\\\"orderInvoiceOpenMark\\\":2,\\\"adjustIsExists\\\":false,\\\"adjustCount\\\":0,\\\"orderFinanceOrgCode\\\":706,\\\"isJDGetCash\\\":true,\\\"adjustId\\\":0,\\\"orderJingdouMoney\\\":0,\\\"ts\\\":\\\"2016-10-11 10:18:11\\\",\\\"buyerCityName\\\":\\\"武汉市\\\",\\\"buyerCountryName\\\":\\\"江汉区\\\",\\\"buyerCoordType\\\":2,\\\"buyerLat\\\":30.6216,\\\"buyerLng\\\":114.244,\\\"buyerIp\\\":1001139541,\\\"orderBuyerRemark\\\":\\\"所购商品如遇缺货，您需要（门店默认配置）：未缺货商品继续配送（缺货商品退款）\\\",\\\"businessTag\\\":\\\"dj_new_cashier;dj_aging_immediately;dj_bld;picking_up;\\\",\\\"equipmentId\\\":\\\"B8F4F80C-EBAD-4223-B666-1784E4F77F37\\\",\\\"buyerPoi\\\":\\\"金雅公寓\\\",\\\"product\\\":[{\\\"skuId\\\":2001716488,\\\"skuName\\\":\\\"潘婷乳液修护优惠装700ml/瓶\\\",\\\"skuIdIsv\\\":\\\"6903148204078\\\",\\\"skuSpuId\\\":0,\\\"skuJdPrice\\\":4900,\\\"skuCount\\\":1,\\\"skuStockOwner\\\":0,\\\"isGift\\\":false,\\\"adjustMode\\\":0,\\\"upcCode\\\":\\\"6903148204078\\\",\\\"categoryId\\\":\\\"21232,21233,22984\\\",\\\"skuStorePrice\\\":4900,\\\"promotionType\\\":1,\\\"promotionId\\\":17640,\\\"skuWeight\\\":0.7979999780654907,\\\"canteenMoney\\\":0},{\\\"skuId\\\":2001716594,\\\"skuName\\\":\\\"康师傅爆椒牛肉面98g/袋\\\",\\\"skuIdIsv\\\":\\\"6920734737901\\\",\\\"skuSpuId\\\":0,\\\"skuJdPrice\\\":300,\\\"skuCount\\\":2,\\\"skuStockOwner\\\":0,\\\"isGift\\\":false,\\\"adjustMode\\\":0,\\\"upcCode\\\":\\\"6920734737901\\\",\\\"categoryId\\\":\\\"20392,20949,20970\\\",\\\"skuStorePrice\\\":300,\\\"promotionType\\\":1,\\\"promotionId\\\":17640,\\\"skuWeight\\\":0.09799999743700027,\\\"canteenMoney\\\":0},{\\\"skuId\\\":2001716721,\\\"skuName\\\":\\\"统一来一桶老坛酸菜牛肉味面120g/桶\\\",\\\"skuIdIsv\\\":\\\"6925303773106\\\",\\\"skuSpuId\\\":0,\\\"skuJdPrice\\\":400,\\\"skuCount\\\":2,\\\"skuStockOwner\\\":0,\\\"isGift\\\":false,\\\"adjustMode\\\":0,\\\"upcCode\\\":\\\"6925303773106\\\",\\\"categoryId\\\":\\\"20392,20949,20952\\\",\\\"skuStorePrice\\\":400,\\\"promotionType\\\":1,\\\"promotionId\\\":17640,\\\"skuWeight\\\":0.11999999731779099,\\\"canteenMoney\\\":0},{\\\"skuId\\\":2001716726,\\\"skuName\\\":\\\"统一100老坛酸菜牛肉面121g/袋\\\",\\\"skuIdIsv\\\":\\\"6925303773915\\\",\\\"skuSpuId\\\":0,\\\"skuJdPrice\\\":250,\\\"skuCount\\\":4,\\\"skuStockOwner\\\":0,\\\"isGift\\\":false,\\\"adjustMode\\\":0,\\\"upcCode\\\":\\\"6925303773915\\\",\\\"categoryId\\\":\\\"20392,20949,20970\\\",\\\"skuStorePrice\\\":250,\\\"promotionType\\\":1,\\\"promotionId\\\":17640,\\\"skuWeight\\\":0.12099999934434891,\\\"canteenMoney\\\":0}],\\\"discount\\\":[{\\\"skuIds\\\":\\\"2001716594,2001716721,2001716488,2001716726\\\",\\\"discountType\\\":4,\\\"discountDetailType\\\":1,\\\"discountCode\\\":\\\"17640\\\",\\\"discountPrice\\\":500}],\\\"appVersion\\\":\\\"iOS3.4.0\\\",\\\"yn\\\":false,\\\"payChannel\\\":9002,\\\"isDeleted\\\":false,\\\"orderGoodsMoney\\\":7300,\\\"orderStockOwner\\\":3,\\\"orderSkuType\\\":0,\\\"isGroupon\\\":false,\\\"orderBaseFreightMoney\\\":200,\\\"orderLadderFreightMoney\\\":0,\\\"orderAdditionFreightMoney\\\":0,\\\"orderBaseReceivableFreight\\\":0,\\\"orderAcceptTime\\\":\\\"2016-10-11 09:51:00\\\",\\\"businessType\\\":1}],\\\"totalPage\\\":1,\\\"page\\\":1}\",\"detail\":\"\",\"code\":\"0\",\"msg\":\"操作成功\"}","success":true}`

	res, err := jsonconv.NewFromString(raw)
	if err != nil {
		log.Error("Failed to parse: %s", err.Error())
		return
	}

	res_str, _ := res.Marshal()
	log.Info("reparsed:\n%s", res_str)
	return
}


func TestJsonMerge() {
	var str_to, str_fr string

	test_func := func(str_to, str_from string) {
		json_to, _ := jsonconv.NewFromString(str_to)
		json_fr, _ := jsonconv.NewFromString(str_fr)

		json_to.MergeFrom(json_fr)
		json_out, _ := json_to.Marshal()
		log.Info("Orig: %s", str_to)
		log.Info("Merg: %s", str_fr)
		log.Info("ResA: %s", json_out)

		json_to, _ = jsonconv.NewFromString(str_to)
		json_to.MergeFrom(json_fr, jsonconv.Option{OverrideArray: true})
		json_out, _ = json_to.Marshal()
		log.Info("ResB: %s", json_out)

		json_to, _ = jsonconv.NewFromString(str_to)
		json_to.MergeFrom(json_fr, jsonconv.Option{OverrideObject: true})
		json_out, _ = json_to.Marshal()
		log.Info("ResC: %s", json_out)

		return
	}

	str_to = `{"string":"strTo","int":0,"obj":{"obj_str":"obj_string","obj_array":[true,"1",2]}}`
	str_fr = `{"string":"strFr","int":1,"obj":{"obj_str":"obj_str_new","obj_array":[false,"10",20],"obj_int":1234}}`
	test_func(str_to, str_fr)

	str_to = `{"string": "orig"}`
	str_fr = `{"string": [1, 2, 3, 4]}`
	test_func(str_to, str_fr)

	str_to = `{"arr": [1, 2 ,3 ,4]}`
	str_fr = `{"arr": ["5", "6", "7", "8"]}`
	test_func(str_to, str_fr)

	str_to = `{"different_type": null}`
	str_fr = `{"different_type": 123.45}`
	test_func(str_to, str_fr)

	str_to = `{"type_1": 123.45, "type_2": "hello"}`
	str_fr = `{"type_1": {"new":"obj"}, "type_2": null}`
	test_func(str_to, str_fr)

	str_to = `{"obj": {"k1": 1, "k2": "2", "k3": true}}`
	str_fr = `{"obj": {"k2": 2, "k3": false, "k4":4.44, "k5": 5.00}}`
	test_func(str_to, str_fr)
	return
}
