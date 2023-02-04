import { SNSClient, PublishCommand } from "@aws-sdk/client-sns"
const DOCKER = process.env.AWS_LAMBDA_FUNCTION_NAME

export const handler = async (event, context) => {
  if (DOCKER) {
    if (DOCKER == 'test_function') { // local docker
    } else { // cloud docker
    }
  }

  const client = new SNSClient({ region: "us-east-1" })
  const command = new PublishCommand({
    Message: 'hello',
    Subject: 'sub',
    TopicArn: 'arn:aws:sns:us-east-1:919759177803:emailer',
  })
  
  let data

  try {
    data = await client.send(command)
  } catch (err) {
    console.log(err)
  } finally {
    return {
      statusCode: 200,
      body: JSON.stringify({
        data
      }),
    }
  }
}

if (!DOCKER) handler()