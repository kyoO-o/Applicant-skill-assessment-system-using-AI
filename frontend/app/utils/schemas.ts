import * as z from "zod";

export const emailSchema = z
  .string()
  .min(1, "Э-мэйл шаардлагатай")
  .email("Зөв э-мэйл хаяг оруулна уу");

export const passwordSchema = z
  .string()
  .min(8, "Нууц үг багадаа 8 тэмдэгттэй байна");

export const requiredPasswordSchema = z.string().min(1, "Нууц үг шаардлагатай");

export const verificationCodeSchema = z
  .string()
  .length(6, "6 оронтой код оруулна уу")
  .regex(/^\d+$/, "Зөвхөн тоо оруулна уу");

export const loginSchema = z.object({
  email: emailSchema,
  password: requiredPasswordSchema,
});

export const registerBaseSchema = z
  .object({
    firstName: z.string().min(1, "Нэр шаардлагатай"),
    lastName: z.string().min(1, "Овог шаардлагатай"),
    email: emailSchema,
    password: passwordSchema,
    confirmPassword: z.string().min(1, "Нууц үг давтах шаардлагатай"),
  })
  .superRefine((data, ctx) => {
    if (data.password !== data.confirmPassword) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Нууц үг таарахгүй байна",
        path: ["confirmPassword"],
      });
    }
  });

export const registerRecruiterSchema = registerBaseSchema.and(
  z.object({
    companyName: z.string().optional(),
    recruiterPosition: z.string().min(1, "Албан тушаал шаардлагатай"),
  }),
);

export const forgotPasswordEmailSchema = z.object({
  email: emailSchema,
});

export const forgotPasswordCodeSchema = z.object({
  code: verificationCodeSchema,
});

export const resetPasswordSchema = z
  .object({
    newPassword: passwordSchema,
    confirmPassword: z.string().min(1, "Нууц үг давтах шаардлагатай"),
  })
  .superRefine((data, ctx) => {
    if (data.newPassword !== data.confirmPassword) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: "Нууц үг таарахгүй байна",
        path: ["confirmPassword"],
      });
    }
  });

export const verifyEmailSchema = z.object({
  code: verificationCodeSchema,
});

export const userProfileSchema = z.object({
  firstName: z.string().min(1, "Нэр шаардлагатай"),
  lastName: z.string().optional(),
});

export const taskSchema = z.object({
  title: z.string().min(1, "Гарчиг шаардлагатай"),
  description: z.string().optional(),
  duration_days: z.coerce.number().int().min(1).optional(),
});

export const jobSchema = z.object({
  title: z.string().min(1, "Ажлын нэр шаардлагатай"),
  contact_info: z.string().min(1, "Холбоо барих мэдээлэл шаардлагатай"),
  type: z.string().min(1, "Ажлын төрөл шаардлагатай"),
  level: z.string().min(1, "Ажлын түвшин шаардлагатай"),
  department: z.string().min(1, "Салбар шаардлагатай"),
  city: z.string().optional(),
  min_salary: z.coerce.number().min(0).optional(),
  max_salary: z.coerce.number().min(0).optional(),
});
