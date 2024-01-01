// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package woker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"

	"github.com/CocaineCong/tangseng/idl/pb/mapreduce"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/CocaineCong/tangseng/types"
)

// MIT H6.824 lab1

func mapper(ctx context.Context, task *mapreduce.MapReduceTask, mapf func(string, string) []*types.KeyValue) {
	// 从文件名读取content
	content, err := os.ReadFile(task.Input)
	if err != nil {
		log.LogrusObj.Error("mapper", err)
		return
	}
	// 将content交给mapf，缓存结果
	intermediates := mapf(task.Input, string(content))

	// 缓存后的结果会写到本地磁盘，并切成R份
	// 切分方式是根据key做hash
	buffer := make([][]*types.KeyValue, task.NReducer)
	for _, intermediate := range intermediates {
		slot := ihash(intermediate.Key) % task.NReducer
		buffer[slot] = append(buffer[slot], intermediate)
	}
	mapOutput := make([]string, 0)
	for i := 0; i < int(task.NReducer); i++ {
		mapOutput = append(mapOutput, writeToLocalFile(int(task.TaskNumber), i, &buffer[i]))
	}
	// R个文件的位置发送给master
	task.Intermediates = mapOutput
	_, err = TaskCompleted(ctx, task)
	if err != nil {
		log.LogrusObj.Errorf("TaskCompleted failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		log.LogrusObj.Errorf("stack trace: \n%+v\n", err)
	}
}

func writeToLocalFile(x int, y int, kvs *[]*types.KeyValue) string {
	dir, _ := os.Getwd()
	tempFile, err := os.CreateTemp(dir, "mr-tmp-*")
	if err != nil {
		fmt.Println(err)
	}
	enc := json.NewEncoder(tempFile)
	for _, kv := range *kvs {
		if err := enc.Encode(&kv); err != nil {
			fmt.Println(err)
		}
	}
	_ = tempFile.Close()
	outputName := fmt.Sprintf("mr-%d-%d", x, y)
	_ = os.Rename(tempFile.Name(), outputName)
	return filepath.Join(dir, outputName)
}
