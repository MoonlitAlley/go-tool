// Copyright (c) 2019 SHUMEI Inc. All rights reserved.

namespace cpp com.shumei.be
namespace php com.shumei.be
namespace java com.shumei.be


// 预测请求
struct PredictRequest {
    1: optional string requestId;  // 唯一标识本次请求
    2: optional string serviceId;  // 唯一标识一个服务
    3: optional string type;       // 请求所属领域类型
    4: optional string organization;  // 唯一标识一个组织
    5: optional string appId;      // 唯一标识一个业务
    6: optional string tokenId;    // 唯一标识一个用户
    7: optional i64 timestamp;     // 客户端时间戳
    8: optional string data;       // 请求数据内容，JSON字符串
    9: optional string eventId;  // 唯一标识一个事件
}

// 风险级别
const string RISK_LEVEL_PASS = "PASS"      // 放行
const string RISK_LEVEL_REVIEW = "REVIEW"  // 需要再次确认
const string RISK_LEVEL_REJECT = "REJECT"  // 组织

// 预测结果
struct PredictResult {
    1: optional i32 score;         // 该服务恒为0
    2: optional string riskLevel;  // 该服务恒为PASS
    3: optional string detail;     // 风险详情，JSON字符串
}

// 预测器异常
exception PredictException {
    1: optional string code;     // 异常代码
    2: optional string message;  // 异常消息
}

// 预测器（服务）
service Predictor {
    PredictResult predict(1:PredictRequest request) throws (1: PredictException e);
}

