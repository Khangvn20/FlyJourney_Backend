package config

import (
    
)

type MomoConfig struct {
    PartnerCode string
    AccessKey   string
    SecretKey   string
    Endpoint    string
    RedirectUrl string
    IpnUrl      string
}

func NewMomoConfig() *MomoConfig {
    return &MomoConfig{
        PartnerCode: "MOMO",
        AccessKey:  "F8BBA842ECF85",
        SecretKey:  "K951B6PE1waDMi640xX08PD3vg6EkVlz",
        Endpoint:    "https://test-payment.momo.vn/v2/gateway/api/create",
        RedirectUrl: "http://localhost:3000/api/v1/payment/success",
        IpnUrl:      "https://554303d5a5c6.ngrok-free.app/api/v1/payment/momo/callback",
    }
}