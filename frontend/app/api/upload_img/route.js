import { NextResponse } from "next/server";
import path from "path";
import { writeFile } from "fs/promises";

export const POST = async (req, res) => {
  const formData = await req.formData();
  const file = formData.get("file");
  if (!file) {
    console.log("filename => ",file)
    return NextResponse.json({ error: "No files received." }, { status: 400 });
  }

  const buffer = Buffer.from(await file.arrayBuffer());
  const filename = Date.now() + file.name.replaceAll(" ", "_");
  try {
    await writeFile(
      path.join(process.cwd(), "/public/uploads/" + filename),
      buffer
    );
    return NextResponse.json({ Message: filename, status: 201 });
  } catch (error) {
    return NextResponse.json({ Message: "Failed", status: 500 });
  }
};
