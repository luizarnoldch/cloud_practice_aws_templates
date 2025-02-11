#!/bin/bash

# Variables
FUNCTION_NAME="HelloWorldFunction"
REGION="us-east-1"
LOG_GROUP_NAME="/aws/lambda/$FUNCTION_NAME"
OUTPUT_FILE="lambda_logs.txt"

# 1. Ejecutar la función Lambda
echo "Ejecutando la función Lambda: $FUNCTION_NAME"
aws lambda invoke \
    --function-name $FUNCTION_NAME \
    --region $REGION \
    output.json

# 2. Obtener los logs de CloudWatch
echo "Obteniendo logs de CloudWatch para el grupo: $LOG_GROUP_NAME"
aws logs filter-log-events \
    --log-group-name $LOG_GROUP_NAME \
    --region $REGION \
    --output text > $OUTPUT_FILE

# 3. Mostrar mensaje de éxito
echo "Logs guardados en $OUTPUT_FILE"
