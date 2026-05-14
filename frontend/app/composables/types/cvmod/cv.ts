export interface WorkExperience {
  company: string;
  position: string;
  start_date: string;
  end_date: string;
  current: boolean;
  description: string;
}

export interface Education {
  school: string;
  degree: string;
  field: string;
  start_date: string;
  end_date: string;
  current: boolean;
  gpa: string;
}

export interface Language {
  name: string;
  level: string;
}

export interface Training {
  name: string;
  organization: string;
  date: string;
  certificate: string;
}

export interface Exam {
  name: string;
  score: string;
  date: string;
}

export interface Internship {
  company: string;
  position: string;
  start_date: string;
  end_date: string;
  description: string;
}

export interface Award {
  name: string;
  organization: string;
  date: string;
  description: string;
}

export interface CVProfile {
  id?: number;
  first_name: string;
  last_name: string;
  date_of_birth: string;
  gender: string;
  national_id: string;
  driver_licenses: string[];
  marital_status: string;
  phone: string;
  email: string;
  address: string;
  about: string;
  work_experiences: WorkExperience[];
  education: Education[];
  personal_skills: string[];
  professional_skills: string[];
  languages: Language[];
  computer_skills: string[];
  art_skills: string[];
  sport_skills: string[];
  trainings: Training[];
  exams: Exam[];
  internships: Internship[];
  awards: Award[];
}

export function emptyCV(): CVProfile {
  return {
    first_name: "",
    last_name: "",
    date_of_birth: "",
    gender: "",
    national_id: "",
    driver_licenses: [],
    marital_status: "",
    phone: "",
    email: "",
    address: "",
    about: "",
    work_experiences: [],
    education: [],
    personal_skills: [],
    professional_skills: [],
    languages: [],
    computer_skills: [],
    art_skills: [],
    sport_skills: [],
    trainings: [],
    exams: [],
    internships: [],
    awards: [],
  };
}

export function emptyWorkExperience(): WorkExperience {
  return { company: "", position: "", start_date: "", end_date: "", current: false, description: "" };
}

export function emptyEducation(): Education {
  return { school: "", degree: "", field: "", start_date: "", end_date: "", current: false, gpa: "" };
}

export function emptyLanguage(): Language {
  return { name: "", level: "" };
}

export function emptyTraining(): Training {
  return { name: "", organization: "", date: "", certificate: "" };
}

export function emptyExam(): Exam {
  return { name: "", score: "", date: "" };
}

export function emptyInternship(): Internship {
  return { company: "", position: "", start_date: "", end_date: "", description: "" };
}

export function emptyAward(): Award {
  return { name: "", organization: "", date: "", description: "" };
}
