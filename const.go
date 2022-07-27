package main

func initConsts() []method {
	escalationMethods := []method{
		method{
			MethodName: "CreateNewPolicyVersion",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:CreatePolicyVersion",
					PolicyBool: true,
				},
			},
		}, method{

			MethodName: "SetExistingDefaultPolicyVersion",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:SetDefaultPolicyVersion",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "CreateEC2WithExistingIP",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PassRole",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "ec2:RunInstances",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "CreateAccessKey",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:CreateAccessKey",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "CreateLoginProfile",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:CreateLoginProfile",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "UpdateLoginProfile",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:UpdateLoginProfile",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "AttachUserPolicy",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:AttachUserPolicy",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "AttachGroupPolicy",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:AttachGroupPolicy",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "AttachRolePolicy",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:AttachRolePolicy",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "sts:AssumeRole",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PutUserPolicy",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PutUserPolicy",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "AddUserToGroup",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:AddUserToGroup",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "UpdateRolePolicyToAssumeIt",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:UpdateAssumeRolePolicy",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "sts:AssumeRole",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PutRolePolicy",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PutRolePolicy",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "sts:AssumeRole",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PassExistingRoleToNewLambdaThenInvoke",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PassRole",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "lambda:CreateFunction",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "lambda:InvokeFunction",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PassExistingRoleToNewLambdaThenTriggerWithNewDynamo",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PassRole",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "lambda:CreateFunction",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "lambda:CreateEventSourceMapping",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "dynamodb:CreateTable",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "dynamodb:PutItem",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PassExistingRoleToNewLambdaThenTriggerWithExistingDynamo",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PassRole",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "lambda:CreateFunction",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "lambda:CreateEventSourceMapping",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PassExistingRoleToNewGlueDevEndpoint",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PassRole",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "glue:CreateDevEndpoint",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "UpdateExistingGlueDevEndpoint",
			PolicySets: []policySet{
				policySet{
					PolicyName: "glue:UpdateDevEndpoint",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PutRolePolicy",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PutRolePolicy",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "sts:AssumeRole",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PassExistingRoleToCloudFormation",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PassRole",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "cloudformation:CreateStack",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "PassExistingRoleToNewDataPipeline",
			PolicySets: []policySet{
				policySet{
					PolicyName: "iam:PassRole",
					PolicyBool: true,
				},
				policySet{
					PolicyName: "datapipeline:CreatePipeline",
					PolicyBool: true,
				},
			},
		}, method{
			MethodName: "EditExistingLambdaFunctionWithRole",
			PolicySets: []policySet{
				policySet{
					PolicyName: "lambda:UpdateFunctionCode",
					PolicyBool: true,
				},
			},
		},
	}
	return escalationMethods
}
