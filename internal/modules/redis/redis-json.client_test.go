package redis

// func TestRedisJson_SetJson(t *testing.T) {
// 	type fields struct {
// 		Client *redis.Client
// 	}
// 	type args struct {
// 		ctx        context.Context
// 		key        string
// 		value      any
// 		expiration time.Duration
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// {
// 		// 	name: "a",
// 		// 	fields: fields{Client: redis.NewClient(&redis.Options{
// 		// 		Addr:       "34.87.154.179:6379",
// 		// 		Username:   "auth",
// 		// 		Password:   "OQwyPmhfrqZWAGq1xHS4RJ5AmE899sA0t2QHV",
// 		// 		DB:         1,
// 		// 		MaxRetries: 1,
// 		// 	})},
// 		// 	args: args{
// 		// 		ctx:        context.Background(),
// 		// 		key:        "test:a",
// 		// 		value:      "test:a",
// 		// 		expiration: 1 * time.Second,
// 		// 	},
// 		// 	wantErr: false,
// 		// },
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			rjs := &JSONClient{
// 				Client: tt.fields.Client,
// 			}
// 			if err := rjs.SetJSON(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration); (err != nil) != tt.wantErr {
// 				t.Errorf("RedisJson.SetJson() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			str := ""
// 			err := rjs.SetJSON(tt.args.ctx, tt.args.key, &str)
// 			if err != nil {
// 				t.Errorf("RedisJson.GetJson() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 			if tt.args.value.(string) != str {
// 				t.Errorf("RedisJson.GetJson() val = %v, got %v", str, tt.args.value.(string))
// 			}

// 			time.Sleep(1 * time.Second)
// 			if err := rjs.GetJSON(tt.args.ctx, tt.args.key, &str); err != redis.Nil {
// 				t.Errorf("RedisJson.GetJson() after 1s got = %v", str)
// 			}
// 		})
// 	}
// }
